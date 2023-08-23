package mysql

import (
	"bookingBackEnd/model"
	"bookingBackEnd/utils"
	"fmt"
	"sync"
)

var onceForBooking sync.Once

type BookingTable struct {
	tableName string
}

var BookingTableInstance *BookingTable

func NewBookingTable() (*BookingTable, error) {
	onceForBooking.Do(func() {
		BookingTableInstance = &BookingTable{
			tableName: utils.ParamsInstance.BookingTableName,
		}
	})
	return BookingTableInstance, nil
}

func (tb *BookingTable) ValidateBooking(info model.InsertBookingInfo) (ret bool, err error) {
	sqlStr := fmt.Sprintf(`
				SELECT COUNT(*) 
				FROM %s
				WHERE room_id = ?
				AND date = ?
				AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))
			`, tb.tableName)
	var count int
	err = DB.Get(&count, sqlStr, info.ClassroomId, info.Date, info.StartTime, info.StartTime, info.EndTime, info.EndTime)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
	}
	ret = count > 0
	return
}

func (tb *BookingTable) InsertBookingInfo(info model.InsertBookingInfo) (err error) {
	sqlStr := fmt.Sprintf(`
				INSERT %s
				(room_id, user_id, date, start_time, end_time)
				VALUES(:room_id, :user_id, :date, :start_time, :end_time)
			`, tb.tableName)
	_, err = DB.NamedExec(sqlStr, &info)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
	}
	return
}

func (tb *BookingTable) DeleteBookingById(id int, userId int) (err error) {
	deleteSqlStr := fmt.Sprintf(`
				DELETE 
				FROM %s
				WHERE id = :id AND user_id = :user_id;
				`, tb.tableName)
	_, err = DB.NamedExec(deleteSqlStr, map[string]interface{}{"id": id, "user_id": userId})
	if err != nil {
		utils.ErrorLogger.Errorf("Error is :%v", err)
		return
	}
	return
}

func (tb *BookingTable) GetBookingList(bookingInfoList *[]model.DetailedBookingInfo, pageNum int, userId int) (err error) {
	sqlStr := fmt.Sprintf(`
				SELECT b.id, a.location, a.floor, a.roomName, a.capacity, a.power, DATE_FORMAT(b.date, "%%Y-%%m-%%d") AS date, b.start_time, b.end_time
				FROM %s AS a
				INNER JOIN %s AS b ON a.id = b.room_id
				WHERE b.user_id = ?
				ORDER BY b.date, b.start_time
				LIMIT ?, ?;
				`, ClassroomTableInstance.tableName, tb.tableName)
	limitCount := utils.ParamsInstance.BookingLimitNumber
	limitOffset := utils.ParamsInstance.BookingLimitNumber * pageNum
	err = DB.Select(bookingInfoList, sqlStr, userId, limitOffset, limitCount)
	return
}

func (tb *BookingTable) GetBookingPeriod(bookingInfoList *[]model.AppointmentAndRoomId, filterClassroomIds string, date string) (err error) {
	sqlStr := fmt.Sprintf(`SELECT room_id, start_time, end_time
	FROM %s
	WHERE date = ? AND room_id in (?);
	`, tb.tableName)
	err = DB.Select(bookingInfoList, sqlStr, date, filterClassroomIds)
	return
}

func (tb *BookingTable) GetBookingStaticsPerFloor(date string) (ret []model.BookingStaticsPerFloor, err error) {
	sqlStr := fmt.Sprintf(`SELECT c.floor,
				COUNT(DISTINCT CASE WHEN a.date = ? THEN a.room_id END) AS bookingCount,
				COUNT(DISTINCT CASE WHEN a.date != ? OR a.date IS NULL THEN c.id END) AS unbookingCount
				FROM %s AS c
				LEFT JOIN %s AS a ON c.id = a.room_id
				GROUP BY c.floor;
	`, ClassroomTableInstance.tableName, tb.tableName)
	err = DB.Select(&ret, sqlStr, date, date)
	return
}

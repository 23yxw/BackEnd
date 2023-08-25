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

func (tb *BookingTable) InsertBookingInfo(info model.InsertBookingInfo) (id int64, err error) {
	sqlStr := fmt.Sprintf(`
				INSERT %s
				(room_id, user_id, date, start_time, end_time)
				VALUES(:room_id, :user_id, :date, :start_time, :end_time)
			`, tb.tableName)
	result, err := DB.NamedExec(sqlStr, &info)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return
	}

	// 获取插入数据的主键ID
	id, err = result.LastInsertId()
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return
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

func (tb *BookingTable) GetBookingList(bookingInfoList *[]model.DetailedBookingInfo, pageNum int, userId int, roomNameRule string) (err error) {
	keyword := "%" + roomNameRule + "%"
	sqlStr := fmt.Sprintf(`
				SELECT b.id, a.location, a.floor, a.roomName, a.capacity, a.power, DATE_FORMAT(b.date, "%%Y-%%m-%%d") AS date, b.start_time, b.end_time, COUNT(*) OVER() AS total_count
				FROM %s AS a
				INNER JOIN %s AS b ON a.id = b.room_id
				WHERE b.user_id = ?
				AND a.roomName LIKE ?
				ORDER BY b.date, b.start_time
				LIMIT ?, ?;
				`, ClassroomTableInstance.tableName, tb.tableName)
	limitCount := utils.ParamsInstance.BookingLimitNumber
	limitOffset := utils.ParamsInstance.BookingLimitNumber * pageNum
	err = DB.Select(bookingInfoList, sqlStr, userId, keyword, limitOffset, limitCount)
	return
}

func (tb *BookingTable) GetBookingPeriod(bookingInfoList *[]model.AppointmentAndRoomId, room_id int, date string) (err error) {
	sqlStr := fmt.Sprintf(`SELECT room_id, start_time, end_time
	FROM %s
	WHERE date = ? AND room_id =?;
	`, tb.tableName)
	err = DB.Select(bookingInfoList, sqlStr, date, room_id)
	return
}

func (tb *BookingTable) GetBookingStaticsPerFloor(date string) (ret []model.BookingStaticsPerFloor, err error) {
	sqlStr := fmt.Sprintf(`SELECT c.floor, COUNT(*) AS bookingCount
				FROM %s AS b
				INNER JOIN %s AS c ON c.id = b.room_id
				WHERE b.date = ?
				GROUP BY c.floor;
	`, tb.tableName, ClassroomTableInstance.tableName)
	err = DB.Select(&ret, sqlStr, date)
	return
}

func (tb *BookingTable) GetPowerStatics(date string) (ret model.PowerStatics, err error) {
	sqlStr := fmt.Sprintf(`SELECT 
				SUM(CASE WHEN c.power = 1 THEN 1 ELSE 0 END) AS powerCount,
				SUM(CASE WHEN c.power = 0 THEN 1 ELSE 0 END) AS noPowerCount
				FROM %s AS b
				JOIN %s AS c ON b.room_id = c.id
				WHERE b.date = ?;
	`, tb.tableName, ClassroomTableInstance.tableName)
	err = DB.Get(&ret, sqlStr, date)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return
	}
	return
}

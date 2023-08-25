package mysql

import (
	"bookingBackEnd/model"
	"bookingBackEnd/utils"
	"fmt"
	"sync"
)

var onceForUserPreference sync.Once

type UserPreferenceTable struct {
	tableName string
}

var UserPreferenceTableInstance *UserPreferenceTable

func NewUserPreferenceTable() (*UserPreferenceTable, error) {
	onceForUserPreference.Do(func() {
		UserPreferenceTableInstance = &UserPreferenceTable{
			tableName: utils.ParamsInstance.UserPreferenceTableName,
		}
	})
	return UserPreferenceTableInstance, nil
}

func (tb *UserPreferenceTable) GetPreferenceClassroomInfo(third_session string, pageNum int) (classroomList []model.DetailedClassroomInfo, err error) {
	// countSqlStr := fmt.Sprintf(`SELECT COUNT(*)
	// 				FROM %s AS up
	// 				INNER JOIN %s AS u ON up.user_id = u.id
	// 				INNER JOIN %s AS c ON up.room_id = c.id
	// 				WHERE u.third_session = ?
	// 				`, tb.tableName, UserTableInstance.tableName, ClassroomTableInstance.tableName)
	// err = DB.Get(&count, countSqlStr, third_session)
	// if err != nil {
	// 	utils.ErrorLogger.Errorf("error is:%v", err)
	// 	return
	// }

	sqlStr := fmt.Sprintf(`SELECT c.id, c.location, c.floor, c.roomName, c.capacity, c.power, c.photo, COUNT(*) OVER() AS total_count
					FROM %s AS up
					INNER JOIN %s AS u ON up.user_id = u.id
					INNER JOIN %s AS c ON up.room_id = c.id
					WHERE u.third_session = ?
					ORDER BY c.id
					LIMIT ?, ?;
					`, tb.tableName, UserTableInstance.tableName, ClassroomTableInstance.tableName)
	limitCount := utils.ParamsInstance.ClassroomImageLimitNumber
	limitOffset := utils.ParamsInstance.ClassroomImageLimitNumber * pageNum
	err = DB.Select(&classroomList, sqlStr, third_session, limitOffset, limitCount)
	return
}

// func (tb *UserPreferenceTable) GetPreferenceUserId(third_session string) (ret []int, err error) {
// 	sqlStr := fmt.Sprintf(`SELECT up.user_id
// 	FROM %s AS up
// 	INNER JOIN %s AS u ON up.user_id = u.id
// 	WHERE u.third_session = ?;
// 	`, tb.tableName, UserTableInstance.tableName)
// 	err = DB.Select(&ret, sqlStr, third_session)
// 	return
// }

func (tb *UserPreferenceTable) InsertUserPreference(userId, roomId int) (err error) {
	sqlStr := fmt.Sprintf(`
				INSERT %s
				(user_id, room_id)
				VALUES(:user_id, :room_id);
			`, tb.tableName)
	_, err = DB.NamedExec(sqlStr, map[string]interface{}{"user_id": userId, "room_id": roomId})
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
	}
	return
}

func (tb *UserPreferenceTable) DeleteUserPreference(roomId int, userId int) (err error) {
	deleteSqlStr := fmt.Sprintf(`
				DELETE 
				FROM %s
				WHERE room_id = :room_id AND user_id = :user_id;
				`, tb.tableName)
	_, err = DB.NamedExec(deleteSqlStr, map[string]interface{}{"room_id": roomId, "user_id": userId})
	if err != nil {
		utils.ErrorLogger.Errorf("Error is :%v", err)
		return
	}
	return
}

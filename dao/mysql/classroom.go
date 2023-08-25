package mysql

import (
	"bookingBackEnd/model"
	"bookingBackEnd/utils"
	"fmt"

	// "strings"
	"sync"
)

var onceForClassroom sync.Once

type ClassroomTable struct {
	tableName string
}

var ClassroomTableInstance *ClassroomTable

func NewClassroomTable() (*ClassroomTable, error) {
	onceForClassroom.Do(func() {
		ClassroomTableInstance = &ClassroomTable{
			tableName: utils.ParamsInstance.ClassroomTableName,
		}
	})
	return ClassroomTableInstance, nil
}

func (tb *ClassroomTable) InsertClassroomInfo(info model.DetailedClassroomInfo) (id int64, err error) {
	sqlStr := fmt.Sprintf(`
				INSERT %s
				(location, floor, roomName, capacity, power, photo)
				VALUES(:location, :floor, :roomName, :capacity, :power, :photo)
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

func (tb *ClassroomTable) GetClassroomInfoById(classroomId int) (ret model.ClassroomAndPhotoInfo, err error) {
	sqlStr := fmt.Sprintf(`
				SELECT id, location, floor, roomName, capacity, power, photo
				FROM %s
				WHERE id = ?;
				`, tb.tableName)
	err = DB.Get(&ret, sqlStr, classroomId)
	if err != nil {
		utils.ErrorLogger.Errorf("Error is :%v", err)
		return
	}
	return
}

func (tb *ClassroomTable) GetClassroomList(classroomList *[]model.CountClassroomInfo, pageNum int, roomNameRule string) (err error) {
	keyword := "%" + roomNameRule + "%"
	sqlStr := fmt.Sprintf(`SELECT id, location, floor, roomName, capacity, power, COUNT(*) OVER() AS total_count
	FROM %s
	WHERE roomName LIKE ?
	ORDER BY id
	LIMIT ?, ?;
	`, tb.tableName)
	limitCount := utils.ParamsInstance.ClassroomListLimitNumber
	limitOffset := utils.ParamsInstance.ClassroomListLimitNumber * pageNum
	err = DB.Select(classroomList, sqlStr, keyword, limitOffset, limitCount)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return
	}
	return
}

func (tb *ClassroomTable) GetDetailedClassroomList(detailedClassroomList *[]model.DetailedClassroomInfo, pageNum int, roomNameRule string) (err error) {
	keyword := "%" + roomNameRule + "%"
	sqlStr := fmt.Sprintf(`SELECT id, location, floor, roomName, capacity, power, photo, COUNT(*) OVER() AS total_count
	FROM %s
	WHERE roomName LIKE ?
	ORDER BY id
	LIMIT ?, ?;
	`, tb.tableName)
	limitCount := utils.ParamsInstance.ClassroomImageLimitNumber
	limitOffset := utils.ParamsInstance.ClassroomImageLimitNumber * pageNum
	err = DB.Select(detailedClassroomList, sqlStr, keyword, limitOffset, limitCount)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return
	}
	return
}

func (tb *ClassroomTable) DeleteClassroomById(classroomId int) (photoUrl string, err error) {
	sqlStr := fmt.Sprintf(`
				SELECT photo 
				FROM %s
				WHERE id = ?;
				`, tb.tableName)
	err = DB.Get(&photoUrl, sqlStr, classroomId)
	if err != nil {
		utils.ErrorLogger.Errorf("Error is :%v", err)
		return
	}

	deleteSqlStr := fmt.Sprintf(`
				DELETE 
				FROM %s
				WHERE id = :id;
				`, tb.tableName)
	_, err = DB.NamedExec(deleteSqlStr, map[string]interface{}{"id": classroomId})
	if err != nil {
		utils.ErrorLogger.Errorf("Error is :%v", err)
		return
	}
	return
}

func (tb *ClassroomTable) UpdateImageById(imageUrl string, id int) (err error) {
	sqlStr := fmt.Sprintf(`UPDATE %s 
					SET photo = :photo
					WHERE id = :id
					`, tb.tableName)
	_, err = DB.NamedExec(sqlStr, map[string]interface{}{"photo": imageUrl, "id": id})
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return
	}
	return
}

func (tb *ClassroomTable) UpdateClassroomInfoById(info model.DetailedClassroomInfo) (err error) {
	// sqlStr := fmt.Sprintf(`UPDATE %s SET `, tb.tableName)
	// i := 1
	// for key := range info {
	// 	if key != "id" {
	// 		sqlStr += (key + "=:" + key)
	// 		if i < len(info) {
	// 			sqlStr += ","
	// 		}
	// 	}
	// 	i += 1
	// }
	// sqlStr += " WHERE id = :id"
	// utils.ErrorLogger.Info(sqlStr)
	sqlStr := fmt.Sprintf(`UPDATE %s 
					SET location = :location, floor = :floor, roomName = :roomName, capacity = :capacity, power = :power, photo = :photo
					WHERE id = :id;
					`, tb.tableName)
	_, err = DB.NamedExec(sqlStr, info)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return
	}
	return
}

func (tb *ClassroomTable) FilterClassroomInfo(floor int, capacity int, power int, pageNum int) (ret []model.DetailedClassroomInfo, err error) {
	sqlStr := fmt.Sprintf(`SELECT id, location, floor, roomName, capacity, power, photo, COUNT(*) OVER() AS total_count
	FROM %s
	WHERE floor = ?
	AND capacity >= ?
	AND power = ?
	ORDER BY id
	LIMIT ?, ?;
	`, tb.tableName)
	limitCount := utils.ParamsInstance.ClassroomImageLimitNumber
	limitOffset := utils.ParamsInstance.ClassroomImageLimitNumber * pageNum
	err = DB.Select(&ret, sqlStr, floor, capacity, power, limitOffset, limitCount)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return
	}
	return
}

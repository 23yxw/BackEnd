package mysql

import (
	"bookingBackEnd/model"
	"bookingBackEnd/utils"
	"fmt"
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

func (tb *ClassroomTable) InsertClassroomInfo(info model.UploadClassroomInfo) (err error) {
	sqlStr := fmt.Sprintf(`
				INSERT %s
				(location, floor, room_name, capacity, power, photo)
				VALUES(:location, :floor, :room_name, :capacity, :power, :photo)
			`, tb.tableName)
	_, err = DB.NamedExec(sqlStr, &info)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
	}
	return
}

func (tb *ClassroomTable) GetClassroomInfoById(classroomId string) (ret model.UpdateClassroomInfo, err error) {
	sqlStr := fmt.Sprintf(`
				SELECT id, location, floor, room_name, capacity, power, photo
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

func (tb *ClassroomTable) GetClassroomList(classroomList *[]model.ClassroomInfo) (err error) {
	sqlStr := fmt.Sprintf(`SELECT id, location, floor, room_name, capacity, power
	FROM %s
	ORDER BY id
	LIMIT ?;
	`, tb.tableName)
	err = DB.Select(classroomList, sqlStr, utils.ParamsInstance.ClassroomListLimitNumber)
	return
}

func (tb *ClassroomTable) GetDetailedClassroomList(detailedClassroomList *[]model.DetailedClassroomInfo) (err error) {
	sqlStr := fmt.Sprintf(`SELECT id, location, floor, room_name, capacity, power, photo
	FROM %s
	ORDER BY id
	LIMIT ?;
	`, tb.tableName)
	err = DB.Select(detailedClassroomList, sqlStr, utils.ParamsInstance.ClassroomImageLimitNumber)
	return
}

func (tb *ClassroomTable) DeleteClassroomById(classroomId string) (photoUrl string, err error) {
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

func (tb *ClassroomTable) UpdateClassroomInfoById(info map[string]interface{}) (err error) {
	sqlStr := fmt.Sprintf(`UPDATE %s SET `, tb.tableName)
	i := 2
	for key := range info {
		if key != "id" {
			sqlStr += (key + "=:" + key)
		}
		if i < len(info) {
			sqlStr += ","
		}
		i += 1
	}
	sqlStr += " WHERE id = :id"

	_, err = DB.NamedExec(sqlStr, info)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
	}
	return
}

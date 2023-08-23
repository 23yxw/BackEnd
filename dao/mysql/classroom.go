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

func (tb *ClassroomTable) InsertClassroomInfo(info model.UploadClassroomInfo) (err error) {
	sqlStr := fmt.Sprintf(`
				INSERT %s
				(location, floor, roomName, capacity, power, photo)
				VALUES(:location, :floor, :roomName, :capacity, :power, :photo)
			`, tb.tableName)
	_, err = DB.NamedExec(sqlStr, &info)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
	}
	return
}

func (tb *ClassroomTable) GetClassroomInfoById(classroomId int) (ret model.UpdateClassroomInfo, err error) {
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

func (tb *ClassroomTable) GetClassroomList(classroomList *[]model.ClassroomInfo, pageNum int) (err error) {
	sqlStr := fmt.Sprintf(`SELECT id, location, floor, roomName, capacity, power
	FROM %s
	ORDER BY id
	LIMIT ?, ?;
	`, tb.tableName)
	limitCount := utils.ParamsInstance.ClassroomListLimitNumber
	limitOffset := utils.ParamsInstance.ClassroomListLimitNumber * pageNum
	err = DB.Select(classroomList, sqlStr, limitOffset, limitCount)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return
	}
	return
}

func (tb *ClassroomTable) GetDetailedClassroomList(detailedClassroomList *[]model.DetailedClassroomInfo, pageNum int) (err error) {
	sqlStr := fmt.Sprintf(`SELECT id, location, floor, roomName, capacity, power, photo
	FROM %s
	ORDER BY id
	LIMIT ?, ?;
	`, tb.tableName)
	limitCount := utils.ParamsInstance.ClassroomImageLimitNumber
	limitOffset := utils.ParamsInstance.ClassroomImageLimitNumber * pageNum
	err = DB.Select(detailedClassroomList, sqlStr, limitOffset, limitCount)
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

func (tb *ClassroomTable) UpdateClassroomInfoById(info map[string]interface{}) (err error) {
	sqlStr := fmt.Sprintf(`UPDATE %s SET `, tb.tableName)
	i := 1
	for key := range info {
		if key != "id" {
			sqlStr += (key + "=:" + key)
			if i < len(info) {
				sqlStr += ","
			}
		}
		i += 1
	}
	sqlStr += " WHERE id = :id"
	utils.ErrorLogger.Info(sqlStr)
	_, err = DB.NamedExec(sqlStr, info)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return
	}
	return
}

func (tb *ClassroomTable) FilterClassroomId(floor int, capacity int, power int) (ret []model.ClassroomInfo, err error) {
	sqlStr := fmt.Sprintf(`SELECT id, location, floor, roomName, capacity, power
	FROM %s
	WHERE floor = ?
	AND capacity >= ?
	AND power = ?;
	`, tb.tableName)
	err = DB.Select(&ret, sqlStr, floor, capacity, power)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return
	}
	return
}

func (tb *ClassroomTable) GetClassroomPowerStatics() (ret model.ClassroomPowerStatics, err error) {
	sqlStr := fmt.Sprintf(`SELECT 
				SUM(CASE WHEN power = 1 THEN 1 ELSE 0 END) AS powerCount,
				SUM(CASE WHEN power = 0 THEN 1 ELSE 0 END) AS noPowerCount
				FROM %s;
	`, tb.tableName)
	err = DB.Get(&ret, sqlStr)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return
	}
	return
}

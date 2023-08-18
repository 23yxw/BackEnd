package mysql

import (
	// "bookingBackEnd/model"
	"bookingBackEnd/utils"
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

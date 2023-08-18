package mysql

import (
	// "bookingBackEnd/model"
	"bookingBackEnd/utils"
	"sync"
)

var onceForUser sync.Once

type UserTable struct {
	tableName string
}

var UserTableInstance *UserTable

func NewUserTable() (*UserTable, error) {
	onceForUser.Do(func() {
		UserTableInstance = &UserTable{
			tableName: utils.ParamsInstance.UserTableName,
		}
	})
	return UserTableInstance, nil
}

package mysql

import (
	// "bookingBackEnd/model"
	"bookingBackEnd/utils"
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

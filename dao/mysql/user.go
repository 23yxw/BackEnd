package mysql

import (
	"bookingBackEnd/model"
	"bookingBackEnd/utils"
	"fmt"
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

func (tb *UserTable) InsertUser(userInfo model.UserInfo) (err error) {
	sqlStr := fmt.Sprintf(`
				INSERT %s
				(third_session, email, password)
				VALUES(:third_session, :email, :password)
			`, tb.tableName)
	_, err = DB.NamedExec(sqlStr, &userInfo)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
	}
	return
}

// GetThirdSession
func (tb *UserTable) GetThirdSession(email string, password string) (ret []model.ThirdSessionInfo, err error) {
	sqlStr := fmt.Sprintf(`
				SELECT third_session, user_type
				FROM %s 
				WHERE email = ? and password = ?;
				`, tb.tableName)

	err = DB.Select(&ret, sqlStr, email, password)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
	}
	return
}

func (tb *UserTable) GetUserIdBythirdsession(third_session string) (ret model.UserIdInfo, err error) {
	sqlStr := fmt.Sprintf(`
				SELECT id, user_type
				FROM %s
				WHERE third_session = ?;
				`, tb.tableName)
	err = DB.Get(&ret, sqlStr, third_session)
	if err != nil {
		utils.ErrorLogger.Errorf("Error is :%v", err)
		return
	}
	return
}

package services

// import (
// 	"bookingBackEnd/dao/mysql"
// 	"bookingBackEnd/model"
// 	"bookingBackEnd/utils"
// )

// func InsertOrUpdateUserPreference(thirdSession string, roomId int) (err error) {
// 	// 查询user_id是否存在
// 	userId, err := mysql.UserPreferenceTableInstance.GetPreferenceUserId(thirdSession)
// 	if err != nil {
// 		utils.ErrorLogger.Errorf("error:%v", err)
// 		return
// 	}

// 	// 插入表项
// 	if len(userId) == 0 {
// 		var userIdInfo model.UserIdInfo
// 		userIdInfo, err = mysql.UserTableInstance.GetUserIdBythirdsession(thirdSession)
// 		if err != nil {
// 			utils.ErrorLogger.Errorf("error:%v", err)
// 			return
// 		}

// 		err = mysql.UserPreferenceTableInstance.InsertUserPreference(userIdInfo.Id, roomId)
// 		if err != nil {
// 			utils.ErrorLogger.Errorf("error:%v", err)
// 			return
// 		}
// 	} else {
// 		err = mysql.UserPreferenceTableInstance.UpdateUserPreference(userId[0], roomId)
// 		if err != nil {
// 			utils.ErrorLogger.Errorf("error:%v", err)
// 			return
// 		}
// 	}
// 	return
// }

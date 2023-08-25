package controller

import (
	"bookingBackEnd/dao/mysql"
	"bookingBackEnd/services"
	"bookingBackEnd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 定义结构体来存储前端传递的数据
type UserPreferenceData struct {
	ThirdSession string `form:"thirdSession" binding:"required"`
	RoomID       int    `form:"roomId" binding:"required"`
}

func InsertUserPreference(c *gin.Context) {
	// 解析前端传递的数据
	var data UserPreferenceData
	if err := c.MustBindWith(&data, binding.Query); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	userIdInfo, err := mysql.UserTableInstance.GetUserIdBythirdsession(data.ThirdSession)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Permission denied", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	err = mysql.UserPreferenceTableInstance.InsertUserPreference(userIdInfo.Id, data.RoomID)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to insert user preference info", "")
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := utils.JsonResponse(0, map[string]interface{}{}, "", "")
	c.JSON(http.StatusOK, ret)
}

func DeleteUserPreference(c *gin.Context) {
	// 解析前端传递的数据
	var data UserPreferenceData
	if err := c.MustBindWith(&data, binding.Query); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	userIdInfo, err := mysql.UserTableInstance.GetUserIdBythirdsession(data.ThirdSession)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Permission denied", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	err = mysql.UserPreferenceTableInstance.DeleteUserPreference(data.RoomID, userIdInfo.Id)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to delete user preference info", "")
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := utils.JsonResponse(0, map[string]interface{}{}, "", "")
	c.JSON(http.StatusOK, ret)
}

func GetPreferenceClassroom(c *gin.Context) {
	var q FilterClassroomListValidateQuery
	if err := c.MustBindWith(&q, binding.Query); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	filterClassroomInfo, err := mysql.UserPreferenceTableInstance.GetPreferenceClassroomInfo(q.ThirdSessionId, q.PageNum)
	if err != nil || len(filterClassroomInfo) == 0 {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have preference classroom", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	results, err := services.GetClassroomData(filterClassroomInfo)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get desired classroom info", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := utils.JsonResponse(0, results, "", "")
	c.JSON(http.StatusOK, ret)
}

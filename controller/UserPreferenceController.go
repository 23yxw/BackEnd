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

func UpsertUserPreference(c *gin.Context) {
	// 解析前端传递的数据
	var data UserPreferenceData
	if err := c.MustBindWith(&data, binding.Query); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	err := services.InsertOrUpdateUserPreference(data.ThirdSession, data.RoomID)
	if err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to insert or update user preference info", "")
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := utils.JsonResponse(0, map[string]interface{}{}, "", "")
	c.JSON(http.StatusOK, ret)
}

func GetPreferenceClassroomAndBookingPeriod(c *gin.Context) {
	thirdSession, exists := c.GetQuery("thirdSessionId")
	if !exists {
		utils.ErrorLogger.Infof("error: no thirdSession param")
		ret := utils.JsonResponse(1, map[string]interface{}{}, "no thirdSession param", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	date, exists := c.GetPostForm("date")
	if !exists {
		utils.ErrorLogger.Infof("error: no date param")
		ret := utils.JsonResponse(1, map[string]interface{}{}, "no date param", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	filterClassroomInfo, err := mysql.UserPreferenceTableInstance.GetPreferenceClassroomInfo(thirdSession)
	if err != nil || len(filterClassroomInfo) == 0 {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have preference classroom", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	results, err := services.FilterClassroomAndTimeSegments(filterClassroomInfo, date)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get desired classroom and time segments", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := utils.JsonResponse(0, results, "", "")
	c.JSON(http.StatusOK, ret)
}

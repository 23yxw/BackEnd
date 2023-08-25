package controller

import (
	"bookingBackEnd/dao/mysql"
	"bookingBackEnd/model"
	"bookingBackEnd/services"
	"bookingBackEnd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type BookingClassroomValidateQuery struct {
	ThirdSessionId string `form:"thirdSessionId" binding:"required"`
	ClassroomID    int    `form:"classroomID" binding:"required"`
}

type DeleteBookingValidateQuery struct {
	ThirdSessionId string `form:"thirdSessionId" binding:"required"`
	Id             int    `form:"id" binding:"required"`
}

type BookingClassroomJson struct {
	Date      string `form:"date" binding:"required"`
	StartTime string `form:"startTime" binding:"required"`
	EndTime   string `form:"endTime" binding:"required"`
}

type BookingListValidateQuery struct {
	ThirdSessionId string `form:"thirdSessionId" binding:"required"`
	PageNum        int    `form:"pageNum"`
	RoomNameRule   string `form:"roomNameRule"`
}

type FilterClassroomListValidateQuery struct {
	ThirdSessionId string `form:"thirdSessionId" binding:"required"`
	PageNum        int    `form:"pageNum"`
}

type BookingFilterQuery struct {
	Floor    int `form:"floor" binding:"required"`
	Capacity int `form:"capacity" binding:"required"`
	Power    int `form:"power"`
	FilterClassroomListValidateQuery
}

type DateValidateQuery struct {
	Date string `form:"date" binding:"required"`
	BookingClassroomValidateQuery
}

func BookingClassroom(c *gin.Context) {
	var q BookingClassroomValidateQuery
	if err := c.MustBindWith(&q, binding.Query); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	var j BookingClassroomJson
	if err := c.MustBindWith(&j, binding.JSON); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Booking time missing", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	isConflict, id, err := services.InsertBookingInfo(q.ThirdSessionId, q.ClassroomID, j.Date, j.StartTime, j.EndTime)
	if err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		if isConflict {
			ret := utils.JsonResponse(2, map[string]interface{}{}, "Invalid booking time", "")
			c.JSON(http.StatusOK, ret)
		} else {
			ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to booking this classroom", "")
			c.JSON(http.StatusOK, ret)
		}
		return
	}

	ret := utils.JsonResponse(0, map[string]interface{}{"id": id}, "", "")
	c.JSON(http.StatusOK, ret)
}

func DeleteBooking(c *gin.Context) {
	var q DeleteBookingValidateQuery
	if err := c.MustBindWith(&q, binding.Query); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	userIdInfo, err := mysql.UserTableInstance.GetUserIdBythirdsession(q.ThirdSessionId)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get user info", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	err = mysql.BookingTableInstance.DeleteBookingById(q.Id, userIdInfo.Id)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to delete booking", "")
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := utils.JsonResponse(0, map[string]interface{}{}, "", "")
	c.JSON(http.StatusOK, ret)
}

func GetBookingList(c *gin.Context) {
	var q BookingListValidateQuery
	if err := c.MustBindWith(&q, binding.Query); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	userIdInfo, err := mysql.UserTableInstance.GetUserIdBythirdsession(q.ThirdSessionId)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Permisssion denied", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	var BookingList []model.DetailedBookingInfo
	err = mysql.BookingTableInstance.GetBookingList(&BookingList, q.PageNum, userIdInfo.Id, q.RoomNameRule)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get classroom info", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := utils.JsonResponse(0, BookingList, "", "")
	c.JSON(http.StatusOK, ret)
}

func FilterClassroomForBooking(c *gin.Context) {
	var q BookingFilterQuery
	if err := c.MustBindWith(&q, binding.Query); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	// var j BookingFilterJson
	// if err := c.MustBindWith(&j, binding.JSON); err != nil {
	// 	utils.ErrorLogger.Errorf("error is: %v", err)
	// 	ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
	// 	c.JSON(http.StatusOK, ret)
	// 	return
	// }

	_, err := mysql.UserTableInstance.GetUserIdBythirdsession(q.ThirdSessionId)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Permisssion denied", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	filterClassroomInfo, err := mysql.ClassroomTableInstance.FilterClassroomInfo(q.Floor, q.Capacity, q.Power, q.PageNum)
	if err != nil || len(filterClassroomInfo) == 0 {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get desired classroom", "")
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

	// results, err := services.FilterClassroomAndTimeSegments(filterClassroomInfo, j.Date)
	// if err != nil {
	// 	utils.ErrorLogger.Errorf("error:%v", err)
	// 	ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get desired classroom and time segments", "")
	// 	c.JSON(http.StatusOK, ret)
	// 	return
	// }

	ret := utils.JsonResponse(0, results, "", "")
	c.JSON(http.StatusOK, ret)
}

func GetBookingPeriodByClassroomId(c *gin.Context) {
	var q DateValidateQuery
	if err := c.MustBindWith(&q, binding.Query); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	_, err := mysql.UserTableInstance.GetUserIdBythirdsession(q.ThirdSessionId)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Permisssion denied", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	results, err := services.GetTimeSegments(q.ClassroomID, q.Date)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get desired classroom and time segments", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := utils.JsonResponse(0, results, "", "")
	c.JSON(http.StatusOK, ret)
}

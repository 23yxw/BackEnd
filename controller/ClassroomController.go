package controller

import (
	"bookingBackEnd/dao/mysql"
	"bookingBackEnd/model"
	"bookingBackEnd/services"
	"bookingBackEnd/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ClassroomInfoJson struct {
	Location string `form:"location" binding:"required"`
	Floor    int    `form:"floor" binding:"required"`
	RoomName string `form:"roomName" binding:"required"`
	Capacity int    `form:"capacity" binding:"required"`
	Power    int    `form:"power" binding:"required"`
}

type UpdateClassroomInfoJson struct {
	Location string `form:"location"`
	Floor    int    `form:"floor"`
	RoomName string `form:"roomName"`
	Capacity int    `form:"capacity"`
	Power    int    `form:"power"`
}

type ClassroomListValidateQuery struct {
	ThirdSessionId string `form:"thirdSessionId" binding:"required"`
	PageNum        int    `form:"pageNum"`
}

type SingleClassroomValidateQuery struct {
	ThirdSessionId string `form:"thirdSessionId" binding:"required"`
	Id             int    `form:"id" binding:"required"`
}

type BookingStaticsValidateQuery struct {
	ThirdSessionId string `form:"thirdSessionId" binding:"required"`
	Date           string `form:"date" binding:"required"`
}

func UploadClassroomInfo(c *gin.Context) {
	thirdSession, exists := c.GetQuery("thirdSessionId")
	if !exists {
		utils.ErrorLogger.Infof("error: no thirdSession param")
		ret := utils.JsonResponse(1, map[string]interface{}{}, "no thirdSession param", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	userIdInfo, err := mysql.UserTableInstance.GetUserIdBythirdsession(thirdSession)
	if err != nil || userIdInfo.UserType != 1 {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Permisssion denied", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	var baseInfo ClassroomInfoJson
	if err := c.ShouldBind(&baseInfo); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	imgUrl, err := GetFilesByType("img", c)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Upload image failed", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	var completeInfo model.UploadClassroomInfo
	completeInfo.Photo = imgUrl
	completeInfo.Capacity = baseInfo.Capacity
	completeInfo.Floor = baseInfo.Floor
	completeInfo.Location = baseInfo.Location
	completeInfo.RoomName = baseInfo.RoomName
	completeInfo.Power = baseInfo.Power

	err = mysql.ClassroomTableInstance.InsertClassroomInfo(completeInfo)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to upload classroom", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := utils.JsonResponse(0, map[string]interface{}{}, "", "")
	c.JSON(http.StatusOK, ret)
}

// func GetClassroomInfo(c *gin.Context) {
// classroomId, exists := c.GetQuery("classroomId")
// if !exists {
// 	ret := utils.JsonResponse(1, map[string]interface{}{}, "classroomId not exists", "")
// 	c.JSON(http.StatusOK, ret)
// 	return
// }
// 	userIdInfo, err := mysql.UserTableInstance.GetUserIdBythirdsession(thirdSession)
// 	if err != nil && userIdInfo.UserType != 1 {
// 		utils.ErrorLogger.Errorf("error:%v", err)
// 		ret := utils.JsonResponse(1, map[string]interface{}{}, "Permisssion denied", "")
// 		c.JSON(http.StatusOK, ret)
// 		return
// 	}

// 	classroomInfo, err := mysql.ClassroomTableInstance.GetClassroomInfoById(classroomId)
// 	if err != nil {
// 		utils.ErrorLogger.Errorf("error:%v", err)
// 		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get classroom info", "")
// 		c.JSON(http.StatusOK, ret)
// 		return
// 	}
// 	classroomInfoMap := utils.StructToMapWithJson(classroomInfo)
// 	ret := utils.JsonResponse(0, classroomInfoMap, "", "")
// 	c.JSON(http.StatusOK, ret)
// }

func GetClassroomList(c *gin.Context) {
	var q ClassroomListValidateQuery
	if err := c.MustBindWith(&q, binding.Query); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	userIdInfo, err := mysql.UserTableInstance.GetUserIdBythirdsession(q.ThirdSessionId)
	if err != nil || userIdInfo.UserType != 1 {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Permisssion denied", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	var classroomList []model.ClassroomInfo
	err = mysql.ClassroomTableInstance.GetClassroomList(&classroomList, q.PageNum)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get classroom info", "")
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := utils.JsonResponse(0, classroomList, "", "")
	c.JSON(http.StatusOK, ret)
}

func GetDetailedClassroomList(c *gin.Context) {
	var q ClassroomListValidateQuery
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

	results, err := services.GetDetailedClassroomList(q.PageNum)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get detailed classroom info", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	ret := utils.JsonResponse(0, results, "", "")
	c.JSON(http.StatusOK, ret)
}

func DeleteClassroom(c *gin.Context) {
	var q SingleClassroomValidateQuery
	if err := c.MustBindWith(&q, binding.Query); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	userIdInfo, err := mysql.UserTableInstance.GetUserIdBythirdsession(q.ThirdSessionId)
	if err != nil || userIdInfo.UserType != 1 {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Permisssion denied", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	photoUrl, err := mysql.ClassroomTableInstance.DeleteClassroomById(q.Id)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to delete classroom", "")
		c.JSON(http.StatusOK, ret)
		return
	}
	os.Remove(photoUrl)
	ret := utils.JsonResponse(0, map[string]interface{}{}, "", "")
	c.JSON(http.StatusOK, ret)
}

func UpdateClassroomInfo(c *gin.Context) {
	var q SingleClassroomValidateQuery
	if err := c.MustBindWith(&q, binding.Query); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	userIdInfo, err := mysql.UserTableInstance.GetUserIdBythirdsession(q.ThirdSessionId)
	if err != nil || userIdInfo.UserType != 1 {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Permisssion denied", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	updateClassroomInfo := make(map[string]interface{})
	updateClassroomInfo["id"] = q.Id

	classroomInfo, err := mysql.ClassroomTableInstance.GetClassroomInfoById(q.Id)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get classroom info", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	// 更新基本信息
	var baseInfo UpdateClassroomInfoJson
	if err := c.ShouldBind(&baseInfo); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough classroom parameters", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	// 将所需更新的非零值元素加入map
	err = utils.StructAddToMap(baseInfo, updateClassroomInfo)
	if err != nil {
		utils.ErrorLogger.Errorf("error: %v", err)
	}

	// 更新照片
	imgUrl, err := GetFilesByType("img", c)
	if err != nil {
		// 前端上传了图像，但是本地保存图像失败
		if imgUrl != "" {
			utils.ErrorLogger.Errorf("error:%v", err)
			ret := utils.JsonResponse(1, map[string]interface{}{}, "Update image failed", "")
			c.JSON(http.StatusOK, ret)
			return
		}
	} else {
		os.Remove(classroomInfo.Photo)
		updateClassroomInfo["photo"] = imgUrl
	}

	err = mysql.ClassroomTableInstance.UpdateClassroomInfoById(updateClassroomInfo)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to update classroom", "")
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := utils.JsonResponse(0, map[string]interface{}{}, "", "")
	c.JSON(http.StatusOK, ret)
}

// func GetClassroomStatics(c *gin.Context) {
// 	var q BookingStaticsValidateQuery
// 	if err := c.ShouldBind(&q); err != nil {
// 		utils.ErrorLogger.Errorf("error is: %v", err)
// 		ret := utils.JsonResponse(1, map[string]interface{}{}, "Do not have enough parameters", "")
// 		c.JSON(http.StatusOK, ret)
// 		return
// 	}

// 	userIdInfo, err := mysql.UserTableInstance.GetUserIdBythirdsession(q.ThirdSessionId)
// 	if err != nil || userIdInfo.UserType != 1 {
// 		utils.ErrorLogger.Errorf("error:%v", err)
// 		ret := utils.JsonResponse(1, map[string]interface{}{}, "Permisssion denied", "")
// 		c.JSON(http.StatusOK, ret)
// 		return
// 	}

// 	powerStatics, err := mysql.ClassroomTableInstance.GetClassroomPowerStatics()
// 	if err != nil {
// 		utils.ErrorLogger.Errorf("error:%v", err)
// 		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get power statics", "")
// 		c.JSON(http.StatusOK, ret)
// 		return
// 	}
// 	powerStaticsMap :=

// 	bookingStatics, err := mysql.BookingTableInstance.GetBookingStaticsPerFloor(q.Date)
// 	if err != nil {
// 		utils.ErrorLogger.Errorf("error:%v", err)
// 		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get booking statics per floor", "")
// 		c.JSON(http.StatusOK, ret)
// 		return
// 	}
// 	bookingStaticsMap := utils.StructArrayToMapWithJson(bookingStatics)

// 	classroomInfoMap := utils.StructToMapWithJson(classroomInfo)
// 	ret := utils.JsonResponse(0, classroomInfoMap, "", "")
// 	c.JSON(http.StatusOK, ret)
// }

func GetFilesByType(fileType string, c *gin.Context) (FilePath string, err error) {
	FileHeader, err := c.FormFile(fileType)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v, no files", err)
		return
	}
	// 构造本地图片存储路径
	filenameBytes := utils.GetMd5()
	file := fmt.Sprintf("%x", filenameBytes)
	FilePath = utils.ParamsInstance.ImageDir + "/" + file + ".jpg"

	// 保存图片到本地，并返回图片路径
	err = c.SaveUploadedFile(FileHeader, FilePath)
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v, save files failed", err)
		return
	}
	return
}

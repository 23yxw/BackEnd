package controller

import (
	"bookingBackEnd/dao/mysql"
	"bookingBackEnd/model"
	"bookingBackEnd/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ClassroomInfoJson struct {
	Location string `form:"location" binding:"required"`
	Floor    string `form:"floor" binding:"required"`
	RoomName string `form:"roomName" binding:"required"`
	Capacity string `form:"capacity" binding:"required"`
	Power    string `form:"power" binding:"required"`
}

type UpdateClassroomInfoJson struct {
	Location string `form:"location"`
	Floor    string `form:"floor"`
	RoomName string `form:"roomName"`
	Capacity string `form:"capacity"`
	Power    string `form:"power"`
}

type ClassroomId struct {
	Id int `form:"id" binding:"required"`
}

func UploadClassroomInfo(c *gin.Context) {
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

func GetClassroomInfo(c *gin.Context) {
	classroomId, exists := c.GetQuery("classroomId")
	if !exists {
		ret := utils.JsonResponse(1, map[string]interface{}{}, "classroomId not exists", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	classroomInfo, err := mysql.ClassroomTableInstance.GetClassroomInfoById(classroomId)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get classroom info", "")
		c.JSON(http.StatusOK, ret)
		return
	}
	classroomInfoMap := utils.StructToMapWithJson(classroomInfo)
	ret := utils.JsonResponse(0, classroomInfoMap, "", "")
	c.JSON(http.StatusOK, ret)
}

func GetClassroomList(c *gin.Context) {
	var classroomList []model.ClassroomInfo
	err := mysql.ClassroomTableInstance.GetClassroomList(&classroomList)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get classroom info", "")
		c.JSON(http.StatusOK, ret)
		return
	}
	ret := utils.JsonResponse(0, classroomList, "", "")
	c.JSON(http.StatusOK, ret)
}

func DeleteClassroom(c *gin.Context) {
	classroomId, exists := c.GetQuery("classroomId")
	if !exists {
		ret := utils.JsonResponse(1, map[string]interface{}{}, "classroomId not exists", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	photoUrl, err := mysql.ClassroomTableInstance.DeleteClassroomById(classroomId)
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
	updateClassroomInfo := make(map[string]interface{})

	classroomId, exists := c.GetQuery("classroomId")
	if !exists {
		ret := utils.JsonResponse(1, map[string]interface{}{}, "classroomId not exists", "")
		c.JSON(http.StatusOK, ret)
		return
	}
	updateClassroomInfo["id"] = classroomId

	classroomInfo, err := mysql.ClassroomTableInstance.GetClassroomInfoById(classroomId)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to get classroom info", "")
		c.JSON(http.StatusOK, ret)
		return
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

	// 更新基本信息
	var baseInfo UpdateClassroomInfoJson
	if err := c.MustBindWith(&baseInfo, binding.JSON); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		ret := utils.JsonResponse(1, map[string]interface{}{}, "Failed to update classroom", "")
		c.JSON(http.StatusOK, ret)
		return
	}

	// 将所需更新的非零值元素加入map
	err = utils.StructAddToMap(baseInfo, updateClassroomInfo)
	if err != nil {
		utils.ErrorLogger.Errorf("error: %v", err)
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

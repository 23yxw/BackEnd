package controller

import (
	"bookingBackEnd/dao/mysql"
	"bookingBackEnd/model"
	"bookingBackEnd/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	uuid "github.com/satori/go.uuid"
)

type LoginValidateJson struct {
	Email    string `form:"email" binding:"required"`
	PassWord string `form:"password" binding:"required"`
}

type RegisterValidateJson struct {
	Email    string `form:"email" binding:"required"`
	PassWord string `form:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var j LoginValidateJson
	if err := c.MustBindWith(&j, binding.JSON); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		return
	}

	ret, err := mysql.UserTableInstance.GetThirdSession(j.Email, j.PassWord)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
	}

	results := make(map[string]interface{})
	if len(ret) == 0 {
		res := utils.JsonResponse(1, results, "Incorrect email or password", "")
		c.JSON(http.StatusOK, res)
	} else {
		results["userType"] = ret[0].UserType
		results["thirdSession"] = ret[0].ThirdSession
		res := utils.JsonResponse(0, results, "", "")
		c.JSON(http.StatusOK, res)
	}
}

func Register(c *gin.Context) {
	var j RegisterValidateJson
	if err := c.MustBindWith(&j, binding.JSON); err != nil {
		utils.ErrorLogger.Errorf("error is: %v", err)
		return
	}

	ret, err := mysql.UserTableInstance.GetThirdSession(j.Email, j.PassWord)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
	}
	results := make(map[string]interface{})
	if len(ret) > 0 {
		// 之前注册过
		res := utils.JsonResponse(1, results, "You have already registered!", "")
		c.JSON(http.StatusOK, res)
	} else {
		thirdSession := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
		var userInfo model.UserInfo
		userInfo.ThirdSession = thirdSession
		userInfo.Email = j.Email
		userInfo.PassWord = j.PassWord
		err := mysql.UserTableInstance.InsertUser(userInfo)
		if err != nil {
			utils.ErrorLogger.Errorf("error:%v", err)
			ret := utils.JsonResponse(1, map[string]interface{}{}, "Registration failed", "")
			c.JSON(http.StatusOK, ret)
			return
		}
		ret := utils.JsonResponse(0, map[string]interface{}{}, "Registration succeed", "")
		c.JSON(http.StatusOK, ret)
	}
}

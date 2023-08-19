package router

import (
	// "bookingBackEnd/services"
	"bookingBackEnd/controller"
	"bookingBackEnd/utils"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	// 初始化router,创建带有日志与恢复中间件的路由
	router := gin.Default()
	// 使用zap logger中间件 and https
	router.Use(utils.GinLogger(), utils.GinRecovery(true))

	user := router.Group("/user")
	{
		// 用户注册或登录
		user.POST("/login", controller.Login)
		user.POST("/register", controller.Register)
		// user.GET("/history")
		// user.POST("/preference")
	}

	classroom := router.Group("/classroom")
	{
		// 获取所有教室信息
		classroom.GET("/list", controller.GetClassroomList)
		// 获取所有教室信息，包括二进制图像内容
		// classroom.GET("/detailedList", controller.GetDetailedClassroomList)
		// 	// 获取单个教室详细信息，包括可预约时间段
		// classroom.GET("/info")
		// 	// 获取某个教室的统计特征
		// 	classroom.GET("/statics")
		classroom.POST("/insert", controller.UploadClassroomInfo)
		classroom.POST("/update", controller.UpdateClassroomInfo)
		classroom.DELETE("/delete", controller.DeleteClassroom)
	}

	// booking := router.Group("/booking")
	// {
	// 	booking.POST("/insert")
	// 	booking.POST("/update")
	// 	booking.DELETE("/delete")
	// }

	return router
}

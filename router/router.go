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
		user.POST("/updatePreference", controller.UpsertUserPreference)
	}

	classroom := router.Group("/classroom")
	{
		// 获取所有教室信息
		classroom.GET("/list", controller.GetClassroomList)
		// 获取所有教室信息，包括图像数据
		classroom.GET("/detailedList", controller.GetDetailedClassroomList)
		// 	// 获取某个教室的统计特征
		// 	classroom.GET("/statics")
		classroom.POST("/insert", controller.UploadClassroomInfo)
		classroom.POST("/update", controller.UpdateClassroomInfo)
		classroom.DELETE("/delete", controller.DeleteClassroom)
	}

	booking := router.Group("/booking")
	{
		booking.POST("/insert", controller.BookingClassroom)
		booking.DELETE("/delete", controller.DeleteBooking)
		booking.GET("/history", controller.GetBookingList)
		// 根据条件获取教室信息，包括可预约时间段
		booking.GET("/info", controller.FilterClassroomAndBookingPeriod)
		booking.GET("/preferenceInfo", controller.GetPreferenceClassroomAndBookingPeriod)
	}

	return router
}

package router

import (
	// "bookingBackEnd/services"
	"bookingBackEnd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	// 初始化router,创建带有日志与恢复中间件的路由
	router := gin.Default()
	// 使用zap logger中间件 and https
	router.Use(utils.GinLogger(), utils.GinRecovery(true))

	wechat := router.Group("/wechat")
	{
		// 用户注册或登录
		wechat.GET("/login", func(c *gin.Context) {
			utils.ErrorLogger.Info("Hello World")
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello, World!",
				"data":    "hello world",
			})
		})
	}

	return router
}

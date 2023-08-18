package main

import (
	"bookingBackEnd/dao/mysql"
	"bookingBackEnd/router"
	"bookingBackEnd/utils"
	"flag"
	"fmt"
	"strconv"
)

func main() {
	// 配置文件路径
	var config_filepath string
	flag.StringVar(&config_filepath, "config", "./conf/conf.ini", "配置文件路径")
	flag.Parse()
	fmt.Println("config_filepath: ", config_filepath)

	// 加载配置文件
	utils.InitConfig(config_filepath)

	// 初始化 logger
	utils.NewLoggerHelper()

	// 加载mysql数据库
	err := mysql.Init()
	if err != nil {
		panic(err)
	}

	// 初始化router
	router := router.NewRouter()
	err = router.Run(":" + strconv.Itoa(utils.ParamsInstance.RouterPort))
	if err != nil {
		utils.ErrorLogger.Error("err is %v", err)
		return
	}
}

/*
 * @Author: RichardoMu
 * @Description:
 * @File: db
 * @Date: 2021/11/18 21:43
 */

package mysql

import (
	"bookingBackEnd/utils"
	"fmt"
	"time"

	//"database/sql"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	DB *sqlx.DB
)

// 任意可序列化数据转为Json,便于查看
func Data2Json(anyData interface{}) string {
	JsonByte, err := json.Marshal(anyData)
	if err != nil {
		utils.ErrorLogger.Errorf("数据序列化为json出错:%s", err.Error())
	}
	return string(JsonByte)
}

// 初始化
func Init() error {
	// "weikai_root:wtDc00sQRQPoF=S#r8D@tcp(qingvoice.mysql.rds.aliyuncs.com:3306)/sight_singing?parseTime=True"
	MysqlParams := utils.ParamsInstance.MysqlParams
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Asia%%2FShanghai", MysqlParams.MysqlUserName, MysqlParams.MysqlUserPwd,
		MysqlParams.MysqlWebPath, MysqlParams.MysqlDatabase)
	//fmt.Println("dns:", dns)
	//dns := "weikai_root:wtDc00sQRQPoF=S#r8D@tcp(qingvoice.mysql.rds.aliyuncs.com:3306)/sight_singing?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
	var err error
	DB, err = sqlx.Open("mysql", dns)
	//defer DB.Close()
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return err
	}
	// 查看是否连接成功
	err = DB.Ping()
	if err != nil {
		utils.ErrorLogger.Errorf("error is:%v", err)
		return err
	}
	DB.SetMaxOpenConns(500)
	DB.SetMaxIdleConns(200)
	DB.SetConnMaxLifetime(10 * time.Minute)

	NewBookingTable()
	NewUserTable()
	NewClassroomTable()
	NewUserPreferenceTable()
	return nil
}

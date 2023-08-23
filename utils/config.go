package utils

import (
	"fmt"
	"os"
	// "strings"

	"gopkg.in/ini.v1"
)

type Params struct {
	MysqlParams      `ini:"mysql"`
	LogDirPathParams `ini:"log_path"`
	RouterParams     `ini:"router"`
	AppParams        `ini:"app"`
	DBTableParams    `ini:"database"`
}

type MysqlParams struct {
	MysqlWebPath  string
	MysqlPort     string
	MysqlUserName string
	MysqlUserPwd  string
	MysqlDatabase string
}

type LogDirPathParams struct {
	ErrorLogPath string
	InfoLogPath  string
}

type RouterParams struct {
	RouterPort int
}

type AppParams struct {
	StartTime                 string
	EndTime                   string
	MinTimeInterval           int
	ImageDir                  string
	ClassroomListLimitNumber  int
	ClassroomImageLimitNumber int
	BookingLimitNumber        int
	MinimumBookingAdvanceTime int
}

type DBTableParams struct {
	BookingTableName        string
	UserTableName           string
	ClassroomTableName      string
	UserPreferenceTableName string
}

var ParamsInstance *Params

/*
go-ini 在字符串中有#号时，需要用``或者"""
https://github.com/go-ini/ini/issues/141
*/

func InitConfig(config_filepath string) {
	cfg, err := ini.Load(config_filepath) //初始化一个cfg
	if err != nil {
		fmt.Printf("load configure file failed, err :%v", err)
		os.Exit(1)
	}
	var params Params
	err = cfg.MapTo(&params)
	if err != nil {
		fmt.Printf("configure file incompatible with `Params` Type, err :%v", err)
		os.Exit(1)
	}
	ParamsInstance = &Params{
		MysqlParams:      params.MysqlParams,
		LogDirPathParams: params.LogDirPathParams,
		RouterParams:     params.RouterParams,
		AppParams:        params.AppParams,
		DBTableParams:    params.DBTableParams,
	}
}

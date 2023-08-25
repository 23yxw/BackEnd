## Environment Requirement
```
Go version 1.20.4
Mysql version 8.0.23
```

## Database setup
```
CREATE DATABASE XX
USE XX
source ./classroom_booking.sql
```

## Project setup
```
go mod tidy // 下载项目需要库
go build -O ./bin/bookingServer main.go //编译
./bin/bookingServer --config ./conf/conf.ini //运行
```

## Parameters
```
// mysql数据库的配置信息，需要改成自己的信息
MysqlWebPath = 172.25.128.1
MysqlPort = 3306
MysqlUserName = root
MysqlUserPwd = `CaF7e871L0GHMSvAah6`
MysqlDatabase = classroom_booking

// 项目的一些参数信息，可以进行自定义修改
StartTime = 08:00:00                //预约开始时间
EndTime = 20:00:00
MinTimeInterval = 1                 //时间段间隔
MinimumBookingAdvanceTime = 2       //预约最短提前时间
ImageDir = ./images
ClassroomListLimitNumber = 4        //分页返回信息的数量
ClassroomImageLimitNumber = 4
BookingLimitNumber = 6
```

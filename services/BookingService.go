package services

import (
	"bookingBackEnd/dao/mysql"
	"bookingBackEnd/model"
	"bookingBackEnd/utils"
	"errors"
	// "fmt"
	// "strings"
	"time"
)

func InsertBookingInfo(thirdSession string, classroomId int, date string, startTime string, endTime string) (isConflict bool, id int64, err error) {
	userIdInfo, err := mysql.UserTableInstance.GetUserIdBythirdsession(thirdSession)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		return
	}

	var bookingInfo model.InsertBookingInfo
	bookingInfo.ClassroomId = classroomId
	bookingInfo.UserId = userIdInfo.Id
	bookingInfo.Date = date
	bookingInfo.StartTime = startTime
	bookingInfo.EndTime = endTime

	// 插入预约记录时判断是否有冲突
	isConflict, err = mysql.BookingTableInstance.ValidateBooking(bookingInfo)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		return
	}
	if isConflict {
		err = errors.New("invalid booking time")
		return
	}

	id, err = mysql.BookingTableInstance.InsertBookingInfo(bookingInfo)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		return
	}
	return
}

// func FilterClassroomAndTimeSegments(filterClassroomInfo []model.ClassroomInfo, date string) (results []map[string]interface{}, err error) {
// 	// 过滤出符合要求的教室
// 	filterClassroomIds := make([]int, len(filterClassroomInfo))
// 	for i, s := range filterClassroomInfo {
// 		filterClassroomIds[i] = s.Id
// 	}

// 	// 将整数数组转换为逗号分隔的字符串
// 	stringArray := make([]string, len(filterClassroomIds))
// 	for i, value := range filterClassroomIds {
// 		stringArray[i] = fmt.Sprint(value)
// 	}
// 	filterClassroomIdString := strings.Join(stringArray, ",")

// 	var bookingInfoList []model.AppointmentAndRoomId
// 	err = mysql.BookingTableInstance.GetBookingPeriod(&bookingInfoList, filterClassroomIdString, date)
// 	if err != nil {
// 		utils.ErrorLogger.Errorf("error:%v", err)
// 		return
// 	}

// 	// 使用map进行分组
// 	groups := make(map[int][]model.AppointmentAndRoomId)
// 	for _, id := range filterClassroomIds {
// 		groups[id] = []model.AppointmentAndRoomId{}
// 		for _, item := range bookingInfoList {
// 			groups[item.RoomID] = append(groups[item.RoomID], item)
// 		}
// 	}

// 	results = make([]map[string]interface{}, 0)
// 	for roomId, bookingInfo := range groups {
// 		var timeSegments []map[string]interface{}
// 		timeSegments, err = getAvailableSlots(utils.ParamsInstance.StartTime, utils.ParamsInstance.EndTime, date, utils.ParamsInstance.MinTimeInterval, bookingInfo)
// 		if err != nil {
// 			utils.ErrorLogger.Errorf("error:%v", err)
// 			return
// 		}
// 		result := make(map[string]interface{})
// 		result["timeSegments"] = timeSegments
// 		for _, s := range filterClassroomInfo {
// 			if roomId == s.Id {
// 				result["classroomInfo"] = s
// 			}
// 		}
// 		results = append(results, result)
// 	}
// 	return
// }

func GetTimeSegments(room_id int, date string) (results []map[string]interface{}, err error) {
	// 取出某教室的预约记录
	var bookingInfoList []model.AppointmentAndRoomId
	err = mysql.BookingTableInstance.GetBookingPeriod(&bookingInfoList, room_id, date)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		return
	}

	// 获取可预约的时间段信息
	results, err = getAvailableSlots(utils.ParamsInstance.StartTime, utils.ParamsInstance.EndTime, date, utils.ParamsInstance.MinTimeInterval, bookingInfoList)
	if err != nil {
		utils.ErrorLogger.Errorf("error:%v", err)
		return
	}
	return
}

func getAvailableSlots(startTime, endTime, date string, interval int, appointments []model.AppointmentAndRoomId) (timeSegments []map[string]interface{}, err error) {
	parserdStartTime, err := time.Parse("15:04:05", startTime)
	if err != nil {
		utils.ErrorLogger.Errorf("error is %v", err)
		return
	}
	parserdEndTime, err := time.Parse("15:04:05", endTime)
	if err != nil {
		utils.ErrorLogger.Errorf("error is %v", err)
		return
	}

	minInterval := time.Duration(interval) * time.Hour

	// 初始化第一个可用时间段
	timeSegments = make([]map[string]interface{}, 0)
	currentTime := parserdStartTime

	for currentTime.Before(parserdEndTime) {
		timeSegment := make(map[string]interface{})
		timeSegment["available"] = 1

		for _, appointment := range appointments {
			parserdAppointmentStartTime, _ := time.Parse("15:04:05", appointment.StartTime)
			parserdAppointmentEndTime, _ := time.Parse("15:04:05", appointment.EndTime)
			if parserdAppointmentStartTime.Before(currentTime.Add(minInterval)) && parserdAppointmentEndTime.After(currentTime) {
				timeSegment["available"] = 0
				break
			}
		}
		if !IsAfterCurrentTime(date, currentTime) {
			timeSegment["available"] = 0
		}
		timeSegment["startTime"] = currentTime.Format("15:04:05")
		timeSegment["endTime"] = currentTime.Add(minInterval).Format("15:04:05")

		timeSegments = append(timeSegments, timeSegment)
		currentTime = currentTime.Add(minInterval)
	}
	return
}

func IsAfterCurrentTime(dateStr string, timeObj time.Time) bool {
	// 解析日期字符串和时间字符串
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.ErrorLogger.Errorf("error is %v", err)
		return false
	}

	// 加载"Asia/Shanghai"时区
	loc, err := time.LoadLocation("Asia/Shanghai") // 北京时间
	// loc, err := time.LoadLocation("Europe/London") // 伦敦时间
	if err != nil {
		utils.ErrorLogger.Errorf("error is %v", err)
		return false
	}

	// 获取当前日期和时间
	currentTime := time.Now().In(loc)
	minimumBookingAdvanceTime := time.Duration(utils.ParamsInstance.MinimumBookingAdvanceTime) * time.Hour

	// 使用解析后的日期和时间组合成完整的时间变量
	combinedTime := time.Date(date.Year(), date.Month(), date.Day(), timeObj.Hour(), timeObj.Minute(), timeObj.Second(), 0, loc)

	// 判断组合后的时间是否晚于当前时间
	return combinedTime.After(currentTime.Add(minimumBookingAdvanceTime))
}

package model

type BaseBookingInfo struct {
	ClassroomId int    `db:"room_id" json:"room_id"`
	Date        string `db:"date" json:"date"`
	StartTime   string `db:"start_time" json:"start_time"`
	EndTime     string `db:"end_time" json:"end_time"`
}

type InsertBookingInfo struct {
	UserId int `db:"user_id" json:"user_id"`
	BaseBookingInfo
}

type DetailedBookingInfo struct {
	ID        int    `db:"id" json:"id"`
	Date      string `db:"date" json:"date"`
	StartTime string `db:"start_time" json:"start_time"`
	EndTime   string `db:"end_time" json:"end_time"`
	BaseInfo
}

type AppointmentAndRoomId struct {
	RoomID    int    `db:"room_id" json:"room_id"`
	StartTime string `db:"start_time" json:"start_time"`
	EndTime   string `db:"end_time" json:"end_time"`
}

type BookingStaticsPerFloor struct {
	Floor          int `db:"floor" json:"floor"`
	BookingCount   int `db:"bookingCount" json:"bookingCount"`
	UnbookingCount int `db:"unbookingCount" json:"unbookingCount"`
}

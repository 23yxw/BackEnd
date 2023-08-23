package model

type BaseInfo struct {
	Location string `db:"location" json:"location"`
	Floor    int    `db:"floor" json:"floor"`
	RoomName string `db:"roomName" json:"roomName"`
	Capacity int    `db:"capacity" json:"capacity"`
	Power    int    `db:"power" json:"power"`
}

type UploadClassroomInfo struct {
	Photo string `db:"photo" json:"photo"`
	BaseInfo
}

type UpdateClassroomInfo struct {
	Id int `db:"id" json:"id"`
	UploadClassroomInfo
}

type ClassroomInfo struct {
	Id int `db:"id" json:"id"`
	BaseInfo
}

type DetailedClassroomInfo struct {
	Photo string `db:"photo" json:"photo"`
	ClassroomInfo
}

type ClassroomPowerStatics struct {
	PowerCount   int `db:"powerCount" json:"powerCount"`
	NoPowerCount int `db:"noPowerCount" json:"noPowerCount"`
}

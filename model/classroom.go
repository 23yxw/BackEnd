package model

type BaseInfo struct {
	Location string `db:"location" json:"location"`
	Floor    int    `db:"floor" json:"floor"`
	RoomName string `db:"roomName" json:"roomName"`
	Capacity int    `db:"capacity" json:"capacity"`
	Power    int    `db:"power" json:"power"`
}

type UpdateClassroomInfo struct {
	Id int `db:"id" json:"id"`
	BaseInfo
}

type ClassroomInfo struct {
	Id int `db:"id" json:"id"`
	BaseInfo
}

type ClassroomAndPhotoInfo struct {
	Photo string `db:"photo" json:"photo"`
	ClassroomInfo
}

type CountClassroomInfo struct {
	TotalCount int `db:"total_count" json:"total_count"`
	ClassroomInfo
}

type DetailedClassroomInfo struct {
	Photo string `db:"photo" json:"photo"`
	CountClassroomInfo
}

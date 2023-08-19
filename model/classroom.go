package model

type BaseInfo struct {
	Location string `db:"location" json:"location"`
	Floor    string `db:"floor" json:"floor"`
	RoomName string `db:"room_name" json:"room_name"`
	Capacity string `db:"capacity" json:"capacity"`
	Power    string `db:"power" json:"power"`
}

type UploadClassroomInfo struct {
	Photo string `db:"photo" json:"photo"`
	BaseInfo
}

type UpdateClassroomInfo struct {
	Id string `db:"id" json:"id"`
	UploadClassroomInfo
}

type ClassroomInfo struct {
	Id string `db:"id" json:"id"`
	BaseInfo
}

type DetailedClassroomInfo struct {
	Id string `db:"id" json:"id"`
	ClassroomInfo
}

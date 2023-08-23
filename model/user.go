package model

type UserInfo struct {
	ThirdSession string `db:"third_session" json:"third_session"`
	Email        string `db:"email" json:"email"`
	PassWord     string `db:"password" json:"password"`
}

type ThirdSessionInfo struct {
	ThirdSession string `db:"third_session" json:"third_session"`
	UserType     int    `db:"user_type" json:"user_type"`
}

type UserIdInfo struct {
	Id       int `db:"id" json:"third_session"`
	UserType int `db:"user_type" json:"user_type"`
}

package model

type User struct {
	Uid      int64  `json:"uid"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

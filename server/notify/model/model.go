package model

type UserMsgID struct {
	UID          int64 `json:"uid"`
	CurrentMsgID int64 `json:"current_msg_id"`
	TotalMsgID   int64 `json:"total_msg_id"`
}

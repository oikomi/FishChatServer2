package model

var (
	HbaseTable = []byte("im")

	HbaseFamilyUser = []byte("user")
	HbaseFamilyMsg  = []byte("msg")

	HbaseColumnSourceUID        = []byte("sourceUID")
	HbaseColumnTargetUID        = []byte("targetUID")
	HbaseColumnGroupID          = []byte("groupID")
	HbaseColumnOnline           = []byte("online")
	HbaseColumnIncrementID      = []byte("incrementID")
	HbaseColumnMsgType          = []byte("msgType")
	HbaseColumnMsgID            = []byte("msgID")
	HbaseColumnMsg              = []byte("msg")
	HbaseColumnAccessServerAddr = []byte("accessServerAddr")
)

type OffsetMsg struct {
	SourceUID int64  `json:"sourceUID"`
	TargetUID int64  `json:"targetUID"`
	MsgID     string `json:"msgID"`
	Msg       string `json:"msg"`
}

type UserMsgID struct {
	UID          int64 `json:"uid"`
	CurrentMsgID int64 `json:"current_msg_id"`
	TotalMsgID   int64 `json:"total_msg_id"`
}

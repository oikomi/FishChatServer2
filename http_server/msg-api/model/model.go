package model

type Login struct {
	Token string `json:"token"`
}

type OfflineMsg struct {
	SourceUID int64
	TargetUID int64
	MsgID     string
	Msg       string
}

type OfflineMsgs struct {
	Msgs []*OfflineMsg
}

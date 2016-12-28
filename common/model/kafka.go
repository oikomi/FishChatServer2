package model

const (
	SendP2PMsgKey   = "send_p2p_msg"
	SendGroupMsgKey = "send_group_msg"
)

type SendP2PMsgKafka struct {
	SourceUID        int64
	TargetUID        int64
	MsgID            string
	Msg              string
	AccessServerAddr string
	Online           bool
}

type SendGroupMsgKafka struct {
	SourceUID int64
	GroupID   int64
	MsgID     string
	Msg       string
}

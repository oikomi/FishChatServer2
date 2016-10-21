package model

const (
	SendP2PMsgKey = "send_p2p_msg"
)

type SendP2PMsgKafka struct {
	UID       int64
	TargetUID int64
	Msg       string
}

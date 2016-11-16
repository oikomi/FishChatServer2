package model

const (
	SendP2PMsgKey = "send_p2p_msg"
)

type SendP2PMsgKafka struct {
	SourceUID int64
	TargetUID int64
	Msg       string
	Online    bool
}

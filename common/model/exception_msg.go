package model

import (
// "github.com/oikomi/FishChatServer2/common/xtime"
)

type ExceptionMsg struct {
	SourceUID int64
	TargetUID int64
	MsgID     string
	Msg       string
}

package model

import (
	"github.com/oikomi/FishChatServer2/common/xtime"
)

type OfflineMsg struct {
	MsgID     int64
	SourceUID int64
	TargetUID int64
	Msg       string
	CTime     xtime.Time
	MTime     xtime.Time
	flag      bool
}

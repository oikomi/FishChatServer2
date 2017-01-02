package model

import (
	"github.com/oikomi/FishChatServer2/common/xtime"
)

type ExceptionMsg struct {
	SourceUID int64
	TargetUID int64
	MsgID     string
	Msg       string
}

type OfflineMsg struct {
	MsgID     string
	SourceUID int64
	TargetUID int64
	Msg       string
	CTime     xtime.Time
	MTime     xtime.Time
	Flag      bool
}

type OfflineMsgs struct {
	Msgs []*OfflineMsg
}

type Group struct {
}

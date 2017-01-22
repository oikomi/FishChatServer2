package dao

import (
	"github.com/oikomi/FishChatServer2/common/dao/xhbase"
)

type HBase struct {
	client *xhbase.Client
}

func NewHBase() *HBase {
	return &HBase{
		client: xhbase.NewClient("127.0.0.1:2181"),
	}
}

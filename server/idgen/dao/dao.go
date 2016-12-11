package dao

import (
// "github.com/garyburd/redigo/redis"
// "github.com/golang/glog"
// "github.com/oikomi/FishChatServer2/common/dao/xredis"
// "github.com/oikomi/FishChatServer2/server/idgen/conf"
// "golang.org/x/net/context"
)

const (
	_keyExceptionMsg = "mge_"
	_keyNormalMsg    = "mgn_"
)

func keyExceptionMsg(msgID string) string {
	return _keyExceptionMsg + msgID
}

func keyNormalMsg(msgID string) string {
	return _keyNormalMsg + msgID
}

type Dao struct {
}

func NewDao() (dao *Dao) {
	return
}

package dao

import (
	// "github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/dao/xredis"
	"github.com/oikomi/FishChatServer2/server/manager/conf"
	"golang.org/x/net/context"
)

const (
	_keyExceptionMsg = "mge_"
)

func keyExceptionMsg(msgID string) string {
	return _keyExceptionMsg + msgID
}

type Dao struct {
	redis *xredis.Pool
}

func NewDao() (dao *Dao) {
	dao = &Dao{
		redis: xredis.NewPool(conf.Conf.Redis.Redis),
	}
	return
}

func (dao *Dao) SetExceptionMsg(ctx context.Context, msgID string, data string) (err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SET", keyExceptionMsg(msgID), data)
	if err != nil {
		glog.Error(err)
	}
	return
}

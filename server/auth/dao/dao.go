package dao

import (
	"github.com/oikomi/FishChatServer2/common/dao/xredis"
	"github.com/oikomi/FishChatServer2/server/auth/conf"
)

type Dao struct {
	redis *xredis.Pool
}

func NewDao() (dao *Dao) {
	dao = &Dao{
		redis: xredis.NewPool(conf.Conf.Redis.Redis),
	}
	return
}

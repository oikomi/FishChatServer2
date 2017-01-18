package dao

import (
	"github.com/oikomi/FishChatServer2/common/dao/xredis"
	"github.com/oikomi/FishChatServer2/server/notify/conf"
)

type Dao struct {
	redis *xredis.Pool
	Mysql *Mysql
}

func NewDao() (dao *Dao) {
	mysql := NewMysql()
	dao = &Dao{
		redis: xredis.NewPool(conf.Conf.Redis.Redis),
		Mysql: mysql,
	}
	return
}

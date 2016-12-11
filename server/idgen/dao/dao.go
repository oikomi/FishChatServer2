package dao

import (
	// "github.com/garyburd/redigo/redis"
	// "github.com/golang/glog"
	// "github.com/oikomi/FishChatServer2/common/dao/xredis"
	"github.com/oikomi/FishChatServer2/server/idgen/conf"
	// "golang.org/x/net/context"
	"github.com/golang/glog"
)

type Dao struct {
	Etcd *Etcd
}

func NewDao() (dao *Dao, err error) {
	e, err := NewEtcd(conf.Conf.Etcd)
	if err != nil {
		glog.Error(err)
	}
	dao = &Dao{
		Etcd: e,
	}
	return
}

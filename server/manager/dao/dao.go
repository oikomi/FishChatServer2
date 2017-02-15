package dao

import (
	"github.com/golang/glog"
)

type Dao struct {
	Redis   *Redis
	MongoDB *MongoDB
	HBase   *HBase
	Mysql   *Mysql
}

func NewDao() (dao *Dao) {
	m, err := NewMongoDB()
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	redis := NewRedis()
	h := NewHBase()
	mysql := NewMysql()
	dao = &Dao{
		Redis:   redis,
		MongoDB: m,
		HBase:   h,
		Mysql:   mysql,
	}
	return
}

package dao

import (
	"github.com/golang/glog"
)

type Dao struct {
	MongoDB *MongoDB
	Kafka   *Kafka
}

func NewDao() (dao *Dao, err error) {
	m, err := NewMongoDB()
	if err != nil {
		glog.Error(err)
		return
	}
	dao = &Dao{
		MongoDB: m,
		Kafka:   NewKafka(),
	}
	return
}

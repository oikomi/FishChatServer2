package dao

import (
	"github.com/golang/glog"
)

type Dao struct {
	MongoDB       *MongoDB
	KafkaProducer *KafkaProducer
	// ES            *ES
}

func NewDao() (dao *Dao, err error) {
	m, err := NewMongoDB()
	if err != nil {
		glog.Error(err)
		return
	}
	KafkaProducer := NewKafkaProducer()
	// es, err := NewES()
	// if err != nil {
	// 	glog.Error(err)
	// 	return
	// }
	dao = &Dao{
		MongoDB:       m,
		KafkaProducer: KafkaProducer,
		// ES:            es,
	}
	go dao.KafkaProducer.HandleError()
	go dao.KafkaProducer.HandleSuccess()
	go dao.KafkaProducer.Process()
	return
}

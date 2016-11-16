package dao

import ()

type Dao struct {
	KafkaProducer *KafkaProducer
}

func NewDao() (dao *Dao, err error) {
	KafkaProducer := NewKafkaProducer()
	dao = &Dao{
		KafkaProducer: KafkaProducer,
	}
	go dao.KafkaProducer.HandleError()
	go dao.KafkaProducer.HandleSuccess()
	go dao.KafkaProducer.Process()
	return
}

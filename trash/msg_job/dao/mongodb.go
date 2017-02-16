package dao

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/dao/mongodb"
	commmodel "github.com/oikomi/FishChatServer2/common/model"
	"github.com/oikomi/FishChatServer2/jobs/msg_job/conf"
)

type MongoDB struct {
	m *mongodb.MongoDB
}

func NewMongoDB() (mdb *MongoDB, err error) {
	m, err := mongodb.NewMongoDB(conf.Conf.MongoDB.MongoDB)
	if err != nil {
		glog.Error(err)
	}
	mdb = &MongoDB{
		m: m,
	}
	return
}

func (m *MongoDB) StoreOfflineMsg(msg *commmodel.OfflineMsg) (err error) {
	c := m.m.Session.DB(conf.Conf.MongoDB.DB).C(conf.Conf.MongoDB.OfflineMsgCollection)
	if err = c.Insert(msg); err != nil {
		glog.Error(err)
	}
	return
}

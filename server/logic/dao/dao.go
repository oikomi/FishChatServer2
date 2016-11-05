package dao

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/dao/mongodb"
	"github.com/oikomi/FishChatServer2/server/logic/conf"
	"github.com/oikomi/FishChatServer2/server/logic/model"
)

// const

type Dao struct {
	m *mongodb.MongoDB
}

func NewDao() (dao *Dao, err error) {
	m, err := mongodb.NewMongoDB(conf.Conf.MongoDB.MongoDB)
	dao = &Dao{
		m: m,
	}
	return
}

func (dao *Dao) StoreOfflineMsg(msg *model.OfflineMsg) (err error) {
	c := dao.m.Session.DB(conf.Conf.MongoDB.DB).C(conf.Conf.MongoDB.OfflineMsgCollection)
	if err = c.Insert(msg); err != nil {
		glog.Error(err)
	}
	return
}

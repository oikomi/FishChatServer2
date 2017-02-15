package dao

// import (
// 	"github.com/golang/glog"
// 	"github.com/oikomi/FishChatServer2/common/dao/mongodb"
// 	commmodel "github.com/oikomi/FishChatServer2/common/model"
// 	"github.com/oikomi/FishChatServer2/server/manager/conf"
// 	"gopkg.in/mgo.v2/bson"
// )

// type MongoDB struct {
// 	m *mongodb.MongoDB
// }

// func NewMongoDB() (mdb *MongoDB, err error) {
// 	m, err := mongodb.NewMongoDB(conf.Conf.MongoDB.MongoDB)
// 	if err != nil {
// 		glog.Error(err)
// 	}
// 	mdb = &MongoDB{
// 		m: m,
// 	}
// 	return
// }

// func (m *MongoDB) GetOfflineMsg(uid int64) (res []*commmodel.OfflineMsg, err error) {
// 	c := m.m.Session.DB(conf.Conf.MongoDB.DB).C(conf.Conf.MongoDB.OfflineMsgCollection)
// 	res = make([]*commmodel.OfflineMsg, 0)
// 	if err = c.Find(bson.M{"targetuid": uid}).All(&res); err != nil {
// 		glog.Error(err)
// 	}
// 	glog.Info(res)
// 	return
// }

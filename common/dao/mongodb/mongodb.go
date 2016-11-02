package mongodb

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/conf"
	"gopkg.in/mgo.v2"
	"time"
)

type MongoDB struct {
	Session *mgo.Session
}

func NewMongoDB(c *conf.MongoDB) (m *MongoDB, err error) {
	session, err := mgo.DialWithTimeout(c.Addrs, time.Duration(c.DialTimeout))
	if err != nil {
		glog.Error(err)
		return
	}
	m = &MongoDB{
		Session: session,
	}
	return
}

package dao

import (
	"github.com/golang/glog"
)

type Dao struct {
	Redis   *Redis
	MongoDB *MongoDB
	HBase   *HBase
}

func NewDao() (dao *Dao) {
	m, err := NewMongoDB()
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	redis := NewRedis()
	h := NewHBase()
	dao = &Dao{
		Redis:   redis,
		MongoDB: m,
		HBase:   h,
	}
	return
}

// func (dao *Dao) SetExceptionMsg(ctx context.Context, msgID string, data string) (err error) {
// 	conn := dao.redis.Get(ctx)
// 	defer conn.Close()
// 	_, err = conn.Do("SET", keyExceptionMsg(msgID), data)
// 	if err != nil {
// 		glog.Error(err)
// 	}
// 	return
// }

// func (dao *Dao) ExceptionMsg(ctx context.Context, msgID string) (res string, err error) {
// 	conn := dao.redis.Get(ctx)
// 	defer conn.Close()
// 	res, err = redis.String(conn.Do("GET", keyExceptionMsg(msgID)))
// 	if err != nil {
// 		glog.Error(err)
// 	}
// 	return
// }

// func (dao *Dao) SetNormalMsg(ctx context.Context, msgID string, data string) (err error) {
// 	conn := dao.redis.Get(ctx)
// 	defer conn.Close()
// 	_, err = conn.Do("SET", keyNormalMsg(msgID), data)
// 	if err != nil {
// 		glog.Error(err)
// 	}
// 	return
// }

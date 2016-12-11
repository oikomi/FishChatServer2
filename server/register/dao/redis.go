package dao

import (
	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/register/conf"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

const (
	_keyOnline = "rgo_"
	_keyAccess = "rga_"
	_keyToken  = "rgt_"
	_online    = 1
	_offline   = 0
)

func keyOnline(uid int64) string {
	return _keyOnline + strconv.FormatInt(uid, 10)
}

func keyAccess(uid int64) string {
	return _keyAccess + strconv.FormatInt(uid, 10)
}

func keyToken(uid int64) string {
	return _keyToken + strconv.FormatInt(uid, 10)
}

func (dao *Dao) Token(ctx context.Context, uid int64) (res string, err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	res, err = redis.String(conn.Do("GET", keyToken(uid)))
	if err != nil {
		glog.Error(err)
	}
	return
}

func (dao *Dao) SetToken(ctx context.Context, uid int64, token string) (err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SETEX", keyToken(uid), int(time.Duration(conf.Conf.Redis.Expire)/time.Second), token)
	if err != nil {
		glog.Error(err)
	}
	return
}

func (dao *Dao) RegisterAccess(ctx context.Context, uid int64, accessAddr string) (err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SET", keyAccess(uid), accessAddr)
	if err != nil {
		glog.Error(err)
	}
	return
}

func (dao *Dao) RouterAccess(ctx context.Context, uid int64) (res string, err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	res, err = redis.String(conn.Do("GET", keyAccess(uid)))
	if err != nil {
		glog.Error(err)
	}
	return
}

func (dao *Dao) GetOnline(ctx context.Context, uid int64) (res int, err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	res, err = redis.Int(conn.Do("GET", keyOnline(uid)))
	if err != nil {
		res = _offline
		glog.Error(err)
	}
	return
}

func (dao *Dao) SetOnline(ctx context.Context, uid int64) (err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SETEX", keyOnline(uid), int(time.Duration(conf.Conf.Redis.Expire)/time.Second), _online)
	if err != nil {
		glog.Error(err)
	}
	return
}

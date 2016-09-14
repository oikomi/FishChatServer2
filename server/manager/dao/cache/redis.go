package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/oikomi/FishChatServer2/server/manager/conf"
)

type RedisCache struct {
	redisPool *redis.Pool
}

func NewRedisCache() (redisCache *RedisCache) {
	cnop := redis.DialConnectTimeout(conf.Conf.Redis.DialTimeout)
	rdop := redis.DialReadTimeout(conf.Conf.Redis.ReadTimeout)
	wdop := redis.DialWriteTimeout(conf.Conf.Redis.WriteTimeout)
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial(conf.Conf.Redis.Proto, conf.Conf.Redis.Addr, cnop, rdop, wdop)
	}, conf.Conf.Redis.Idle)
	redisCache = &RedisCache{
		redisPool: redisPool,
	}
	return
}

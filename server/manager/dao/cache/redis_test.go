package cache

import (
	"testing"
)

func TestRedis(t *testing.T) {
	redisCache := NewRedisCache()
	if redisCache == nil {
		t.Error("failed")
	}
}

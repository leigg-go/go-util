package redis

import (
	"fmt"
	"github.com/bmizerany/assert"
	"github.com/go-redis/redis"
	"github.com/leigg-go/go-util/_redis"
	"testing"
	"time"
)

var opts = &redis.Options{
	Addr:        "192.168.40.131:63790",
	Password:    "111111",
	DB:          0,
	PoolSize:    3,
	MaxRetries:  3,
	IdleTimeout: 60 * time.Second,
}

func TestMustInitDefClient(t *testing.T) {

	_redis.MustInitDefClient(opts)
	assert.Panic(t, fmt.Errorf("_redis: DefClient already exists"), func() {
		_redis.MustInitDefClient(opts)
	})
	defer _redis.Close()

	_, err := _redis.DefClient.Set("k", "v", 1*time.Second).Result()
	assert.Equal(t, err, nil)

	s, _ := _redis.DefClient.Get("k").Result()
	assert.Equal(t, s, "v")

	time.Sleep(1 * time.Second)

	s, err = _redis.DefClient.Get("k").Result()
	assert.Equal(t, true, _redis.IsNilErr(err))
}

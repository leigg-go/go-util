package redis

import (
	"fmt"
	"github.com/bmizerany/assert"
	"github.com/go-redis/redis"
	"github.com/leigg-go/go-util/_redis"
	"log"
	"testing"
	"time"
)

var opts = &redis.Options{
	Addr:        "127.0.0.1:6379",
	Password:    "123",
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

func TestHScan(t *testing.T) {
	_redis.MustInitDefClient(opts)
	defer _redis.Close()
	hk := "hk1"
	//defer _redis.DefClient.Del(hk)

	for i := 0; i < 1000; i++ {
		err := _redis.DefClient.HSet(hk, fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i)).Err()
		if err != nil {
			t.Errorf("%v", err)
		}
	}

	var (
		s      []string
		cursor uint64
		err    error
	)

	for {
		s, cursor, err = _redis.DefClient.HScan(hk, cursor, "", 10).Result()
		if _redis.IsExecErr(err) {
			t.Errorf("%v", err)
			break
		}
		fmt.Printf("kv-pairs: %d\n", len(_redis.StringMap(s)))

		if cursor == 0 {
			log.Printf("2 %v", err)
			break
		}
	}
}

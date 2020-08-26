package lock

import (
	"github.com/go-redis/redis"
	"github.com/leigg-go/go-util/_lock"
	"github.com/leigg-go/go-util/_redis"
	"github.com/leigg-go/go-util/_util"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
	"time"
)

var opts = &redis.Options{
	Addr:         "127.0.0.1:6379",
	Password:     "",
	DB:           0,
	DialTimeout:  2 * time.Second,
	ReadTimeout:  3 * time.Second,
	WriteTimeout: 3 * time.Second,
	MinIdleConns: 1,
	IdleTimeout:  3 * time.Second,
}

var globalCounter int

// N个goroutine并发，但最终只有一个得到锁，并执行 globalCounter++
var goroutineCount = 500

// Lock共享redis-key
const DistributedLockInRedisKey = "DistributedLockInRedis"

// expire值应该略大于task执行耗时
const DistributedLockInRedisExpire = time.Second * 5

func NoCurrencyTask(t *testing.T, i int, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	lock := _lock.NewDistributedLockByRedis(_redis.DefClient, DistributedLockInRedisKey, nil, DistributedLockInRedisExpire)
	// 获取锁操作需设置超时，因为是通过网络, 注意这个超时不要设置太小(需大于正常操作耗时)
	err := lock.Lock(time.Millisecond * 50)
	if err != nil {
		assert.Equal(t, err, _lock.TimeoutErr)
		log.Printf("goroutine %d lock timeout!", i)
		return
	}
	log.Printf("goroutine %d got the rights!", i)
	// 模拟耗时，届时只有一个协程会获得锁，其他均超时
	time.Sleep(time.Second * 1)

	globalCounter++
	err = lock.UnLock()
	_util.PanicIfErr(err, nil)
}

func TestNewDistributedLockInRedis(t *testing.T) {
	_redis.MustInitDefClient(opts)
	defer _redis.Close()

	var wg sync.WaitGroup
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go NoCurrencyTask(t, i, &wg)
	}
	wg.Wait()
	assert.Equal(t, globalCounter, 1)

	// 上面的所有goroutine应该释放锁，下面再执行一次，应该可以获得锁
	wg.Add(1)
	NoCurrencyTask(t, 1000, nil)
	assert.Equal(t, globalCounter, 2)
}

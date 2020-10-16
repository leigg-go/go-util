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
	Password:     "123",
	DB:           0,
	DialTimeout:  2 * time.Second,
	ReadTimeout:  3 * time.Second,
	WriteTimeout: 3 * time.Second,
	MinIdleConns: 1,
	IdleTimeout:  3 * time.Second,
}

var globalCounter int

// N个goroutine并发，但最终只有一个得到锁，并执行 globalCounter++
var goroutineCount = 100

// Lock共享redis-key
const DistributedLockInRedisKey = "DistributedLockInRedis"

// expire值应该略大于task执行耗时
const DistributedLockInRedisExpire = time.Second * 5

func noCurrencyTask(t *testing.T, i int, wg *sync.WaitGroup, lock _lock.DistributedLock) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	// 每个goroutine应该独享一个lock对象，但是lock key是相同的
	_, err := lock.Lock()
	if err != nil {
		assert.Equal(t, err, _lock.TimeoutErr)
		log.Printf("goroutine %d lock timeout!", i)
		return
	}
	log.Printf("goroutine %d got the lock!", i)
	// 模拟耗时，届时只有一个协程会获得锁，其他均超时
	time.Sleep(time.Second * 1)

	globalCounter++
	err = lock.UnLock()
	_util.PanicIfErr(err, nil)
}

func TestNewDistributedLockInRedis(t *testing.T) {
	_redis.MustInitDefClient(opts)
	defer _redis.Close()

	// 获取锁操作需设置超时，>估计网络耗时
	// 另外还需设置重试，如果不重试，lock只会执行一次redis SetNX操作
	opt := _lock.LockOption{Timeout: time.Millisecond * 50, Retry: true}

	lock := _lock.NewDistributedLockByRedis(_redis.DefClient, DistributedLockInRedisKey, nil, DistributedLockInRedisExpire, opt)

	var wg sync.WaitGroup
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		// 不同的协程应该使用一个同key不同对象的lock，一个协程只能释放自己加的锁
		go noCurrencyTask(t, i, &wg, lock.Copy())
	}
	wg.Wait()
	assert.Equal(t, globalCounter, 1)

	// 上面的所有goroutine应该释放锁，下面再执行一次，应该可以获得锁
	noCurrencyTask(t, goroutineCount+1, nil, lock.Copy())
	assert.Equal(t, globalCounter, 2)
}

func TestLockErr(t *testing.T) {
	_redis.MustInitDefClient(opts)
	defer _redis.Close()

	opt := _lock.LockOption{Timeout: time.Millisecond * 50, Retry: true}

	// 每个goroutine应该独享一个lock对象，但是lock key是相同的
	lock := _lock.NewDistributedLockByRedis(_redis.DefClient, DistributedLockInRedisKey, nil, DistributedLockInRedisExpire, opt)
	// 获取锁操作需设置超时，因为是通过网络, 注意这个超时不要设置太小(需大于正常操作耗时)
	_, err := lock.Lock()
	if err != nil {
		t.Fatalf("Should be nil, but got err:%v", err)
	}
	// 上面还没释放锁，同一个内存对象再去获取锁，就会panic
	_, err = lock.Lock()
	assert.Equal(t, _lock.AlreadyAcquiredErr, err)

	err = lock.UnLock()
	if err != nil {
		t.Fatalf("Should be nil, but got err:%v", err)
	}
	err = lock.UnLock()
	assert.Equal(t, _lock.NotLockedErr, err)
}

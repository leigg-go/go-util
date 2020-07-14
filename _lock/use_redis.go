package _lock

import (
	"errors"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

// redis实现的分布式锁

type DistributedLocker interface {
	Lock(timeout time.Duration) error
	UnLock() error
}

type DistributedLockInRedis struct {
	redisC *redis.Client
	key    string
	value  []byte
	expire time.Duration

	mutex  sync.Mutex
	locked bool
}

func NewDistributedLockInRedis(c *redis.Client, k string, v []byte, expire time.Duration) DistributedLocker {
	return &DistributedLockInRedis{redisC: c, key: k, value: v, expire: expire}
}

var (
	TimeoutErr   = errors.New("acquire lock timeout")
	NotLockedErr = errors.New("acquire lock failed")
)

// Lock try to acquire `lock`
func (l *DistributedLockInRedis) Lock(timeout time.Duration) (err error) {
	pollIntvl := time.NewTicker(time.Millisecond * 50)
	deadline := time.Now().Add(timeout)

	for {
		if time.Now().Sub(deadline) >= 0 {
			return TimeoutErr
		}

		r := l.redisC.SetNX(l.key, l.value, l.expire)
		if r.Err() != nil {
			return r.Err()
		}
		ok := r.Val()
		if ok {
			l.mutex.Lock()
			l.locked = true
			l.mutex.Unlock()
			return
		}
		<-pollIntvl.C
	}
}

// Lock try to release `lock`
func (l *DistributedLockInRedis) UnLock() (err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.locked {
		return l.redisC.Del(l.key).Err()
	}
	return NotLockedErr
}

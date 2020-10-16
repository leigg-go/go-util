package _lock

import (
	"errors"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

// Distributed Lock in redis.
type DistributedLockInRedis struct {
	redisC *redis.Client
	key    string
	value  []byte
	expire time.Duration
	// mutex allow to share a Lock in multi goroutines
	mu sync.RWMutex
	// `isLocked` maintains lock state in multi goroutines, it avoids unlocking lock
	// that do not belong to current lock creator
	isLocked bool

	// lock() args

	// if Retry is true, Lock() would continued to Retry until lockTimeout is reached
	retry bool
	// Lock method use lockTimeout as calling timeout arg, it should be a slightly
	// great than network time consuming for one time redis action.
	// default is 1 second
	lockTimeout time.Duration
}

func NewDistributedLockByRedis(c *redis.Client, k string, v []byte, expire time.Duration, opt ...LockOption) DistributedLock {
	l := &DistributedLockInRedis{redisC: c, key: k, value: v, expire: expire}
	if len(opt) > 0 {
		o := opt[0]
		l.retry = o.Retry
		l.lockTimeout = o.Timeout
	}
	if l.lockTimeout <= 0 {
		l.lockTimeout = time.Second
	}
	return l
}

var (
	TimeoutErr         = errors.New("DistributedLockInRedis: acquire lock timeout")
	AlreadyAcquiredErr = errors.New("DistributedLockInRedis: already acquired lock")
	NotLockedErr       = errors.New("DistributedLockInRedis: acquire lock failed")
)

// Lock try to acquire `lock`. If the lock has been acquired, it will return AlreadyAcquiredErr.
// If Retry[0] is true, it would be return immediately when setNX failed(not err), otherwise,
// there would be continued Retry until timeout.
func (l *DistributedLockInRedis) Lock() (ok bool, err error) {
	if l.Locked() {
		err = AlreadyAcquiredErr
		return
	}
	pollIntvl := time.NewTicker(time.Millisecond * 10)
	defer pollIntvl.Stop()

	deadline := time.Now().Add(l.lockTimeout)
	for {
		if time.Now().Sub(deadline) >= 0 {
			err = TimeoutErr
			return
		}

		r := l.redisC.SetNX(l.key, l.value, l.expire)
		if r.Err() != nil {
			err = r.Err()
			return
		}
		if r.Val() {
			l.setLockState(true)
			ok = true
			return
		}
		if !l.retry {
			return
		}

		<-pollIntvl.C
	}
}

// Lock try to release `lock`. If the lock has not been acquired, it will return NotLockedErr
func (l *DistributedLockInRedis) UnLock() error {
	if l.Locked() {
		if err := l.redisC.Del(l.key).Err(); err != nil {
			return err
		}
		l.setLockState(false)
		return nil
	}
	return NotLockedErr
}

func (l *DistributedLockInRedis) Locked() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.isLocked
}

func (l *DistributedLockInRedis) setLockState(state bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.isLocked = state
}

// Copy copy all attrs but mutex and lock state
func (l *DistributedLockInRedis) Copy() DistributedLock {
	return &DistributedLockInRedis{
		redisC:      l.redisC,
		key:         l.key,
		value:       l.value,
		expire:      l.expire,
		mu:          sync.RWMutex{},
		isLocked:    false,
		retry:       l.retry,
		lockTimeout: l.lockTimeout,
	}
}

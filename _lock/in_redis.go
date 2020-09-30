package _lock

import (
	"errors"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

// Distributed Lock in redis.

type DistributedLock interface {
	Lock(timeout time.Duration) error
	UnLock() error
	Copy() DistributedLock
}

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
}

func NewDistributedLockByRedis(c *redis.Client, k string, v []byte, expire time.Duration) DistributedLock {
	return &DistributedLockInRedis{redisC: c, key: k, value: v, expire: expire}
}

var (
	TimeoutErr         = errors.New("DistributedLockInRedis: acquire lock timeout")
	AlreadyAcquiredErr = errors.New("DistributedLockInRedis: already acquired lock")
	NotLockedErr       = errors.New("DistributedLockInRedis: acquire lock failed")
)

// Lock try to acquire `lock`. If the lock has been acquired, it will return AlreadyAcquiredErr
func (l *DistributedLockInRedis) Lock(timeout time.Duration) error {
	if l.Locked() {
		return AlreadyAcquiredErr
	}
	pollIntvl := time.NewTicker(time.Millisecond * 10)
	defer pollIntvl.Stop()
	deadline := time.Now().Add(timeout)

	for {
		if time.Now().Sub(deadline) >= 0 {
			return TimeoutErr
		}

		r := l.redisC.SetNX(l.key, l.value, l.expire)
		if r.Err() != nil {
			return r.Err()
		}
		if r.Val() {
			l.setLockState(true)
			return nil
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
		redisC:   l.redisC,
		key:      l.key,
		value:    l.value,
		expire:   l.expire,
		mu:       sync.RWMutex{},
		isLocked: false,
	}
}

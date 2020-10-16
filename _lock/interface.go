package _lock

import "time"

type DistributedLock interface {
	Lock() (bool, error)
	UnLock() error
	Copy() DistributedLock
}

type LockOption struct {
	Retry   bool
	Timeout time.Duration
}

package _go

import (
	"fmt"
	"sync"
)

// AsyncTask用于异步启动N个协程，并可以在当前协程与其同步
type SafeAsyncTask struct {
	tasks []func()
	wg    sync.WaitGroup
	mutex sync.RWMutex
	err   error // capture first err happened in these goroutines
}

func NewSafeAsyncTask() *SafeAsyncTask {
	return &SafeAsyncTask{wg: sync.WaitGroup{}}
}

func (a *SafeAsyncTask) AddTask(f func()) {
	a.tasks = append(a.tasks, f)
}

func (a *SafeAsyncTask) schedule() {
	for _, f := range a.tasks {
		a.wg.Add(1)
		go func(f func()) {
			defer func() {
				if err := recover(); err != nil {
					a.SetErr(fmt.Errorf("%v", err))
				}
				a.wg.Done()
			}()
			f()
		}(f)
	}
}

func (a *SafeAsyncTask) RunAndWait() {
	a.schedule()
	a.wg.Wait()
}

func (a *SafeAsyncTask) Clear() {
	a.tasks = nil
}

func (a *SafeAsyncTask) SetErr(err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.err == nil {
		a.err = err
	}
}

func (a *SafeAsyncTask) Err() error {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.err
}

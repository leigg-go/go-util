package _go

// AsyncTask用于异步启动一个协程，并可以在当前协程与其同步
type AsyncTask struct {
	done  chan struct{}
	Error error
}

func NewAsyncTask() *AsyncTask {
	return &AsyncTask{done: make(chan struct{})}
}

func (a *AsyncTask) Run(f func()) {
	go func() {
		f()
		a.done <- struct{}{}
	}()
}

func (a *AsyncTask) Wait() {
	<-a.done
}

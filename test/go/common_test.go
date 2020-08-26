package _go

import (
	"github.com/leigg-go/go-util/_go"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestAsyncTask(t *testing.T) {
	s := &http.Server{Addr: ":10241"}

	ak := _go.NewSafeAsyncTask()
	ak.AddTask(func() {
		log.Printf("server1 started...")
		//panic("111")
		ak.SetErr(s.ListenAndServe())

	})
	ak.AddTask(func() {
		log.Printf("stop task started...")
		time.Sleep(time.Second)
		_ = s.Close()
	})

	ak.RunAndWait()

	log.Println("server closed, err:", ak.Err()) // http.ErrServerClosed
}

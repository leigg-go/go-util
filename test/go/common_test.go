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

	ak := _go.NewAsyncTask()
	ak.Run(func() {
		log.Printf("server started...")
		ak.Error = s.ListenAndServe()
	})

	time.Sleep(time.Second)
	s.Close()

	ak.Wait()
	log.Println("server closed, err:", ak.Error)
}

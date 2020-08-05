package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/leigg-go/go-util/_redis"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

/*
go-redis库测试

NOTE: 只有类get命令在key不存在时err会等于redis.Nil, 其他如del命令在key不存在时err是nil。
*/
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

func initClient() {
	if _redis.DefClient == nil || _redis.DefClient.Ping().Err() != nil {
		_redis.DefClient = nil
		_redis.MustInitDefClient(opts)
	}
}

func TestMustInitDefClient(t *testing.T) {
	_redis.MustInitDefClient(opts)
	assert.Panics(t, func() { _redis.MustInitDefClient(opts) }, fmt.Errorf("_redis: DefClient already exists"))
	defer _redis.Close()

	_, err := _redis.DefClient.Set("k", "v", 1*time.Second).Result()
	assert.Equal(t, err, nil)

	s, _ := _redis.DefClient.Get("k").Result()
	assert.Equal(t, s, "v")

	time.Sleep(1 * time.Second)

	s, err = _redis.DefClient.Get("k").Result()
	assert.Equal(t, true, _redis.IsNilErr(err))
}

func TestHScan(t *testing.T) {
	initClient()
	defer _redis.Close()
	hk := "hk1"
	//defer _redis.DefClient.Del(hk)

	for i := 0; i < 1000; i++ {
		err := _redis.DefClient.HSet(hk, fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i)).Err()
		if err != nil {
			t.Errorf("%v", err)
		}
	}

	var (
		s      []string
		cursor uint64
		err    error
	)

	for {
		s, cursor, err = _redis.DefClient.HScan(hk, cursor, "", 10).Result()
		if _redis.IsExecErr(err) {
			t.Errorf("%v", err)
			break
		}
		fmt.Printf("kv-pairs: %d\n", len(_redis.StringMap(s)))

		if cursor == 0 {
			log.Printf("2 %v", err)
			break
		}
	}
}

func TestInfoClients(t *testing.T) {
	initClient()
	defer _redis.Close()

	for i := 0; i < 3; i++ {
		//r := _redis.DefClient.ClientID()  // redis 5.0 命令
		r, err := _redis.DefClient.Do("info", "clients").String()
		assert.Equal(t, err, nil)
		fmt.Printf("目前的clients: %s\n", r)
		time.Sleep(time.Second * 2)
	}
}

func TestBasicCmds(t *testing.T) {
	initClient()
	defer _redis.Close()

	err1 := _redis.DefClient.SetNX("TestBasicCmds", "x", time.Second).Val()
	err2 := _redis.DefClient.SetNX("TestBasicCmds", "x", time.Second).Val()
	log.Printf("err1:%v err2:%v\n", err1, err2)

	// hash test
	_redis.DefClient.HSet("TestBasicCmds_hash", "k1", "123")
	r := _redis.DefClient.HGet("TestBasicCmds_hash", "k2")
	log.Println(r.Result())        // 123 <nil>
	log.Println("Err():", r.Err()) // if field does not exist, err=redis.Nil
	log.Println(r.Val())           // 123

	delNotFoundKeyErr := _redis.DefClient.HDel("TestBasicCmds_hash_not_found", "k").Err()
	log.Printf("delNotFoundKeyErr:%v", delNotFoundKeyErr) // nil
}

package _redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/leigg-go/go-util/_util"
	"sync"
)

/*
Redis registry, app could init redis client by call them directly.
*/

var (
	lock      sync.Mutex
	DefClient *redis.Client
)

// init default pool on registry
func MustInitDefClient(opts *redis.Options) {
	lock.Lock()
	defer lock.Unlock()
	_util.Must(DefClient == nil, fmt.Errorf("_redis: DefClient already exists"))

	DefClient = newClient(opts)
}

func newClient(opts *redis.Options) *redis.Client {
	cli := redis.NewClient(opts)
	_, err := cli.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("_redis: %v", err))
	}
	return cli
}

func Close() error {
	if DefClient == nil {
		return nil
	}
	return DefClient.Close()
}

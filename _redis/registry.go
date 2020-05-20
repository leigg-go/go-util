package _redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/leigg-go/go-util/_util"
	"sync"
)

/*
Redis registry, project could init redis client by call them directly.
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

	DefClient = newCli(opts)
}

func newCli(opts *redis.Options) *redis.Client {
	cli := redis.NewClient(opts)
	_, err := cli.Ping().Result()
	if err != nil {
		panic(err)
	}
	return cli
}

func Close() error {
	return DefClient.Close()
}

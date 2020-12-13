package _redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/leigg-go/go-util/_util"
	"sync"
)

/*
Redis registry, app can directly call them to initialize redis client.
*/

var (
	lock      sync.Mutex
	DefClient *redis.Client
)

// init default pool on registry
func MustInitDef(opts *redis.Options) {
	lock.Lock()
	defer lock.Unlock()
	_util.Must(DefClient == nil, fmt.Errorf("_redis: DefClient already exists"))

	DefClient = newClient(opts)
}

func MustInit(opts *redis.Options) *redis.Client {
	return newClient(opts)
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
	lock.Lock()
	defer lock.Unlock()
	err := DefClient.Close()
	DefClient = nil
	return err
}

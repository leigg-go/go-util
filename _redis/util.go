package _redis

import "github.com/go-redis/redis"

func IsNilErr(err error) bool {
	return err == nil || err == redis.Nil
}

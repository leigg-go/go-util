package _redis

import (
	"github.com/go-redis/redis"
	"github.com/leigg-go/go-util/_typ/_str"
)

func IsNilErr(err error) bool {
	return err == nil || err == redis.Nil
}

func IsExecErr(err error) bool {
	return !IsNilErr(err)
}

func StringMap(s []string) map[string]string {
	var m = make(map[string]string, len(s)/2)
	var tmpK string

	_str.Each(s, func(i int, elem string) {
		if i%2 == 0 {
			tmpK = elem
			return
		}
		m[tmpK] = elem
	})
	return m
}

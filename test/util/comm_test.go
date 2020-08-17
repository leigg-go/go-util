package util

import (
	"fmt"
	"github.com/leigg-go/go-util/_util"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestRandInt(t *testing.T) {
	var min, max = 10, 20
	for i := 0; i < 100000; i++ {
		v := _util.RandInt(min, max)
		log.Println(v)
		assert.True(t, min <= v && max >= v)
	}
}

func TestShortUrl(t *testing.T) {
	u1 := "http://baidu.com"
	newU, err := _util.ShortUrl(u1)
	log.Println(newU, err)
}

func TestFuncName(t *testing.T) {
	// 当前函数全路径： github.com/leigg-go/go-util/test/util.TestFuncName
	fmt.Println(_util.GetRunningFuncName()) // github.com/leigg-go/go-util/test/util.TestFuncName

	// 纯函数名
	fmt.Println(_util.GetFuncName(TestShortUrl, ".")) // TestShortUrl
	// 包含pkg的函数名
	fmt.Println(_util.GetFuncName(TestShortUrl, "/")) // util.TestShortUrl
	// 包含全路径的函数名
	fmt.Println(_util.GetFuncName(TestShortUrl)) // github.com/leigg-go/go-util/test/util.TestShortUrl
}

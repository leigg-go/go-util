package util

import (
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

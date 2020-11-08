package _str

import (
	"fmt"
	"strconv"
)

func MustToInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("_str: %v", err))
	}
	return int(i)
}

func MustToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("_str: %v", err))
	}
	return i
}

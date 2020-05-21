package _str

import (
	"fmt"
	"strconv"
)

func ToInt64(s string, must ...bool) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil && len(must) > 0 && must[0] {
		panic(fmt.Sprintf("_str: %v", err))
	}
	return i, err
}

func MustToInt64(s string) int64 {
	i, _ := ToInt64(s, true)
	return i
}

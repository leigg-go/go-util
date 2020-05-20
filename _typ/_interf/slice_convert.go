package _interf

import (
	"fmt"
	"github.com/leigg-go/go-util/_typ/_err"
	"github.com/leigg-go/go-util/_typ/_float"
	"github.com/leigg-go/go-util/_typ/_int"
	"github.com/leigg-go/go-util/_typ/_str"
	"strings"
)

// 后面需要什么类型继续增加, 但只能加通用的
func ToSliceInterface(slice interface{}) []interface{} {
	typ := fmt.Sprintf("%T", slice)
	if !strings.HasPrefix(typ, "[]") {
		panic(fmt.Sprintf("<%s> not slice type", typ))
	}
	var ret []interface{}
	switch slice.(type) {
	case []int:
		_int.IntEach(slice.([]int), func(seq int, elem int) {
			ret = append(ret, elem)
		})
	case []int16:
		_int.Int16Each(slice.([]int16), func(seq int, elem int16) {
			ret = append(ret, elem)
		})
	case []int32:
		_int.Int32Each(slice.([]int32), func(seq int, elem int32) {
			ret = append(ret, elem)
		})
	case []int64:
		_int.Int64Each(slice.([]int64), func(seq int, elem int64) {
			ret = append(ret, elem)
		})
	case []string:
		_str.StrEach(slice.([]string), func(seq int, elem string) {
			ret = append(ret, elem)
		})
	case []float64:
		_float.Float64Each(slice.([]float64), func(seq int, elem float64) {
			ret = append(ret, elem)
		})
	case []error:
		_err.ErrEach(slice.([]error), func(seq int, elem error) {
			ret = append(ret, elem)
		})
	default:
		panic(fmt.Sprintf("<%s> not support convert to slice", typ))
	}
	return ret
}

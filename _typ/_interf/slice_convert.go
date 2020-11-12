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
// 如果参数不是slice会panic
func ToSliceInterface(slice interface{}) []interface{} {
	typ := fmt.Sprintf("%T", slice)
	if !strings.HasPrefix(typ, "[]") {
		panic(fmt.Sprintf("<%s> is not slice type", typ))
	}
	var ret []interface{}
	switch slice.(type) {
	case []string:
		_str.Each(slice.([]string), func(seq int, elem string) {
			ret = append(ret, elem)
		})
	case []int:
		_int.IntEach(slice.([]int), func(seq int, elem int) {
			ret = append(ret, elem)
		})
	case []int8:
		_int.Int8Each(slice.([]int8), func(seq int, elem int8) {
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
	case []uint:
		_int.UIntEach(slice.([]uint), func(seq int, elem uint) {
			ret = append(ret, elem)
		})
	case []uint8:
		_int.UInt8Each(slice.([]uint8), func(seq int, elem uint8) {
			ret = append(ret, elem)
		})
	case []uint32:
		_int.UInt32Each(slice.([]uint32), func(seq int, elem uint32) {
			ret = append(ret, elem)
		})
	case []uint64:
		_int.UInt64Each(slice.([]uint64), func(seq int, elem uint64) {
			ret = append(ret, elem)
		})
	case []float32:
		_float.Float32Each(slice.([]float32), func(seq int, elem float32) {
			ret = append(ret, elem)
		})
	case []float64:
		_float.Float64Each(slice.([]float64), func(seq int, elem float64) {
			ret = append(ret, elem)
		})
	case []error:
		_err.Each(slice.([]error), func(seq int, elem error) {
			ret = append(ret, elem)
		})
	default:
		panic(fmt.Sprintf("<%s> not support convert to slice", typ))
	}
	return ret
}

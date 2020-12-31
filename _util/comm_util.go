package _util

import (
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func InCollection(elem interface{}, coll []interface{}) bool {
	for _, e := range coll {
		if e == elem {
			return true
		}
	}
	return false
}

func PanicIfErr(err interface{}, ignoreErrs []error, printText ...string) {
	if err != nil {
		for _, e := range ignoreErrs {
			if e == err {
				return
			}
		}
		if len(printText) > 0 {
			panic(printText[0])
		}
		panic(err)
	}
}

func Must(condition bool, err error) {
	if !condition {
		panic(err)
	}
}

func If(condition bool, then func(), _else ...func()) {
	if condition {
		if then != nil {
			then()
		}
	} else {
		for _, f := range _else {
			f()
		}
	}
}

func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

var shanghai, _ = time.LoadLocation("Asia/Shanghai")

const (
	TimeLayoutSec        = "2006-01-02 15:04:05"
	TimeLayoutSecNoSpace = "2006-01-02_15-04-05"
	TimeLayoutDay        = "2006-01-02"
)

func LoadShanghaiTimeFromStr(s string) (time.Time, error) {
	return time.ParseInLocation(TimeLayoutSec, s, shanghai)
}

// 获取指定函数的名称, split:分割符，`.`获取纯函数名， `/`获取带pkg的函数名，如 _util.GetFuncName
func GetFuncName(funcObj interface{}, split ...string) string {
	fn := runtime.FuncForPC(reflect.ValueOf(funcObj).Pointer()).Name()
	if len(split) > 0 {
		fs := strings.Split(fn, split[0])
		return fs[len(fs)-1]
	}
	return fn
}

// 当前运行的函数名, split:分割符，不传就是获取全路径的函数名称
// split 传入 `.`获取纯函数名， `/`获取带pkg的函数名，如 _util.GetRunningFuncName
func GetRunningFuncName(split ...string) string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	fn := runtime.FuncForPC(pc[0]).Name()

	if len(split) > 0 {
		fs := strings.Split(fn, split[0])
		return fs[len(fs)-1]
	}
	return fn
}

// skip=1 为调用者位置，skip=2为调用者往上一层的位置，以此类推
// return-example: /develop/go/test_go/tmp_test.go:88
func FileWithLineNum(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%v:%v", file, line)
}

func IsDBErr(err error) bool {
	if err != nil && err != gorm.ErrRecordNotFound {
		return true
	}
	return false
}

func IsRedisErr(err error) bool {
	if err != nil && err != redis.Nil {
		return true
	}
	return false
}

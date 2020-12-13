package _util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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

func AnyErr(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
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
	TimeLayoutSec = "2006-01-02 15:04:05"
	TimeLayoutDay = "2006-01-02"
)

func LoadShanghaiTimeFromStr(s string) (time.Time, error) {
	return time.ParseInLocation(TimeLayoutSec, s, shanghai)
}

// TODO: complete this method
func shortUrl(oldUrl string) (string, error) {
	servUrl := "https://sina.lt/api.php?from=w&url=%s&site=dwz.date"
	type rsp struct {
		Result string
		Data   interface{}
	}

	bb := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, bb)
	_, _ = encoder.Write([]byte(oldUrl))
	_ = encoder.Close()

	req, _ := http.NewRequest("GET", fmt.Sprintf(servUrl, bb.String()), nil)
	req.Header.Add("referer", "https://sina.lt/")
	req.Header.Add("PHPSESSID", "s7qtpi42pr73u6r73p6q3fj2n2")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.81 Safari/537.36 SE 2.X MetaSr 1.0")

	var cli = http.Client{
		Timeout: time.Second * 5,
	}
	r, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	//body, err := ioutil.ReadAll(r.Body)
	//log.Println(string(body))
	//if err != nil {
	//	return "", err
	//}

	var jsonRsp = new(rsp)
	err = json.NewDecoder(r.Body).Decode(jsonRsp)
	if err != nil {
		return "", err
	}
	if jsonRsp.Result != "ok" {
		log.Println(jsonRsp.Data)
		return "", errors.New(jsonRsp.Result)
	}
	return jsonRsp.Data.(string), nil
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

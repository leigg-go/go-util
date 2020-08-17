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
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"
)

/*
这是一种资源清理模式。首先进程退出的方式有三种：
1. init时某些操作返回err(主协程内的操作)，需要退出
2. 运行时的panic
3. ctrl+C或其他外部信号
=========================
对于第一种情况，在函数返回err时可直接panic，main函数会捕捉并通知OnProcessExit处理
对于第二种，只要不是子协程内发生的panic，main的defer可以捕捉并处理;如果是子协程panic(预期内的),最后通过某种方式通知主协程，而不是直接panic
对于第三种情况就很简单，信号监听和资源清理都在OnProcessExit内部完成
*/
// 监听[进程退出信号]的协程，完成资源释放工作
func OnProcessExit(doWhenClose func()) (chan<- struct{}, <-chan struct{}) {
	done := make(chan struct{})
	shouldExit := make(chan struct{})
	go func() {
		// 监听进程外部信号
		sysSignalChan := make(chan os.Signal)
		signal.Notify(sysSignalChan,
			syscall.SIGHUP, // 终端线挂断
			syscall.SIGINT, // 键盘中断
			//syscall.SIGKILL, // kill信号无法捕捉
			syscall.SIGTERM, // 软件终止
		)
		var onSignal bool
		select {
		case <-shouldExit:
			log.Println("OnProcessExit read chan-shouldExit!")
		case s := <-sysSignalChan:
			onSignal = true
			log.Printf("OnCloseSignal: %s\n", s)
		}
		signal.Stop(sysSignalChan)
		close(sysSignalChan)
		doWhenClose()
		log.Println("OnProcessExit complete!")

		close(done)

		if onSignal {
			os.Exit(0)
		}
		// 不是signal不要退出（否则看不到panic信息），外面调用了ctx.Cancel()，那必须在外面发生panic
	}()
	return shouldExit, done
}

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
		var intErrs []interface{}
		for _, e := range ignoreErrs {
			intErrs = append(intErrs, e)
		}
		if InCollection(err, intErrs) {
			return
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

type SvcWithClose interface {
	Close() error
}

func CloseSvcSafely(manySvc []SvcWithClose) []error {
	var (
		errs []error
		err  error
	)
	for _, s := range manySvc {
		if reflect.ValueOf(s).IsNil() {
			continue
		}
		if err = s.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

var shanghai, _ = time.LoadLocation("Asia/Shanghai")
var simpleLayout = "2006-01-02 15:04:05"

func LoadShanghaiTimeFromStr(s string) (time.Time, error) {
	return time.ParseInLocation(simpleLayout, s, shanghai)
}

func ShortUrl(oldUrl string) (string, error) {
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

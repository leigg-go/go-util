package _util

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

func PanicIfErr(err error, ignoreErrs []error, printText ...string) {
	if err != nil {
		var intErrs []interface{}
		for _, e := range ignoreErrs {
			intErrs = append(intErrs, e)
		}
		if InCollection(err, intErrs) {
			return
		}
		if len(printText) > 0 {
			panic(fmt.Sprintf(printText[0], err))
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

type SvcWithClose interface {
	Close() error
}

func CloseSvcSafely(manySvc []SvcWithClose) []error {
	var (
		errs []error
		err  error
	)
	for _, i := range manySvc {
		if i == nil {
			continue
		}
		if err = i.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

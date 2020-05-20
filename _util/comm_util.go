package _util

import (
	"context"
	"fmt"
	"github.com/leigg-go/go-util/_typ/_interf"
	"log"
	"os"
	"os/signal"
	"syscall"
)

/*
这是一种资源清理模式。首先进程退出的方式有三种：
1. init时某些操作返回err，需要退出
2. 运行时的panic
3. crtl+C或其他外部信号
=========================
对于第一种情况，在函数返回err时可直接panic，main函数会捕捉并通知OnProcessExit处理
对于第二种，只要不是子协程内发生的panic，main的defer可以捕捉并处理
对于第三种流程就很简单，信号监听和资源清理都在OnProcessExit内部完成。
*/
// 监听[进程退出信号]的协程函数，完成资源释放工作
func OnProcessExit(doWhenClose func(), mainProcessCtx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		// 监听进程外部信号
		sysSignalChan := make(chan os.Signal)
		signal.Notify(sysSignalChan,
			syscall.SIGHUP, // 终端线挂断
			syscall.SIGINT, // 键盘中断
			//syscall.SIGKILL, // 杀死程序,这个信号无法捕捉
			syscall.SIGTERM, // 软件终止
		)
		var onSignal bool
		select {
		case <-mainProcessCtx.Done():
			log.Println("OnProcessExit mainProcessCtx.Done!")
		case s := <-sysSignalChan:
			onSignal = true
			log.Printf("OnCloseSignal: %s\n", s)
		}
		signal.Stop(sysSignalChan)
		close(sysSignalChan)
		doWhenClose()
		log.Println("OnProcessExit complete!")

		done <- struct{}{}

		if onSignal {
			os.Exit(0)
		}
		// 不是signal不要退出（否则看不到panic信息），外面调用了ctx.Cancel()，那外面必须发生panic
	}()
	return done
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
		if InCollection(err, _interf.ToSliceInterface(ignoreErrs)) {
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

func IfElse(condition bool, then func(), _else func()) {
	if condition {
		if then != nil {
			then()
		}
	} else {
		if _else != nil {
			_else()
		}
	}
}

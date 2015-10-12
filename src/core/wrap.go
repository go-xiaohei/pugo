package core

import (
	"fmt"
	"gopkg.in/inconshreveable/log15.v2"
	"runtime"
	"sync"
	"time"
)

var (
	wg          WaitGroup
	wgLogFormat = "goroutine.%s"
)

// 调用Wrap的Go程
func Wrap(funcName string, fn func()) {
	wg.Wrap(funcName, fn)
}

func WrapWait() {
	wg.Wait()
}

type WaitGroup struct {
	sync.WaitGroup
}

func (w *WaitGroup) Wrap(funName string, fn func()) {
	w.Add(1)
	go func() {
		t := time.Now()
		fn()
		w.Done()
		log15.Debug(fmt.Sprintf(wgLogFormat, funName), "goroutine", runtime.NumGoroutine(), "duration", time.Since(t).Seconds()*1000)
		// 强制退出
		runtime.Gosched()
	}()
}

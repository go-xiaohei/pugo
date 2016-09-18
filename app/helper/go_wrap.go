package helper

import (
	"runtime"
	"sync"
)

var wg WaitGroup

// GoWrap add global function to global Wait
func GoWrap(funcName string, fn func()) {
	wg.Wrap(funcName, fn)
}

// GoWait wait all wrapping goroutine in global WaitGroup
func GoWait() {
	wg.Wait()
}

// GoGroup define a single goroutine wrap
type GoGroup struct {
	wg   WaitGroup
	name string
}

// NewGoGroup create new group
func NewGoGroup(name string) *GoGroup {
	return &GoGroup{
		name: name,
	}
}

// Wrap add a function to Group
func (g *GoGroup) Wrap(name string, fn func()) {
	g.wg.Wrap(name, fn)
}

// Wait wait all goroutine in this Group
func (g *GoGroup) Wait() {
	g.wg.Wait()
}

// WaitGroup struct
type WaitGroup struct {
	sync.WaitGroup
}

// Wrap add goroutine function wrap
func (w *WaitGroup) Wrap(funName string, fn func()) {
	w.Add(1)
	go func() {
		// t := time.Now()
		fn()
		w.Done()
		// log15.Debug("GoWrap|%s|%.1fms|%d", funName, time.Since(t).Seconds()*1000, runtime.NumGoroutine())
		// make goroutine free
		runtime.Gosched()
	}()
}

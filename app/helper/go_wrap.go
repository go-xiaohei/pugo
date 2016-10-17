package helper

import "sync"

var wg WaitGroup

// GoWrap add global function to global Wait
func GoWrap(funcName string, fn func() error) {
	wg.Wrap(funcName, fn)
}

// GoWait wait all wrapping goroutine in global WaitGroup
func GoWait() {
	wg.Wait()
}

// GoWrapErrors return all errors in this global WaitGroup
func GoWrapErrors() []error {
	return wg.Errors()
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
func (g *GoGroup) Wrap(name string, fn func() error) {
	g.wg.Wrap(name, fn)
}

// Wait wait all goroutine in this Group
func (g *GoGroup) Wait() {
	g.wg.Wait()
}

// Errors return all errors when running goroutines in this group
func (g *GoGroup) Errors() []error {
	return g.wg.Errors()
}

// WaitGroup struct
type WaitGroup struct {
	sync.WaitGroup
	errors []error
}

// Wrap add goroutine function wrap
func (w *WaitGroup) Wrap(funName string, fn func() error) {
	w.Add(1)
	go func() {
		// t := time.Now()
		if err := fn(); err != nil {
			w.errors = append(w.errors, err)
		}
		w.Done()
		// log15.Debug("GoWrap|%s|%.1fms|%d", funName, time.Since(t).Seconds()*1000, runtime.NumGoroutine())
	}()
}

// Errors return errors from all wrap function
func (w *WaitGroup) Errors() []error {
	return w.errors
}

package helper

import (
	"errors"
	"runtime"
	"sync/atomic"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWorker(t *testing.T) {
	var a int64 = 1
	normalFn := func() error {
		atomic.AddInt64(&a, 1)
		return nil
	}
	errorFn := func() error {
		atomic.AddInt64(&a, 1)
		return errors.New("error")
	}

	Convey("Worker", t, func() {
		Convey("AddWorkerFunc", func() {
			w := NewWorker(0)
			So(w.chans, ShouldHaveLength, runtime.NumCPU())
			for i := 0; i < 100; i++ {
				w.AddFunc(normalFn)
			}
			So(w.funcs, ShouldHaveLength, 100)
		})

		Convey("RunOnce", func() {
			w := NewWorker(0)
			for i := 0; i < 100; i++ {
				w.AddFunc(normalFn)
			}
			w.RunOnce()
			So(a, ShouldEqual, 101)
		})

		Convey("WorkerFuncWithError", func() {
			w := NewWorker(0)
			j := 0
			for i := 0; i < 100; i++ {
				if i%2 == 0 {
					j++
					w.AddFunc(errorFn)
				} else {
					w.AddFunc(normalFn)
				}
			}
			So(w.funcs, ShouldHaveLength, 100)
			w.RunOnce()
			So(w.Errors(), ShouldHaveLength, j)
		})
	})
}

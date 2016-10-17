package helper

import (
	"sync/atomic"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var i int32

func TestGlobalWrap(t *testing.T) {
	Convey("GlobalWrap", t, func() {
		GoWrap("+1", func() {
			atomic.AddInt32(&i, 1)
		})
		GoWrap("+2", func() {
			atomic.AddInt32(&i, 2)
		})
		GoWrap("+3", func() {
			atomic.AddInt32(&i, 3)
		})
		GoWait()
		So(i, ShouldEqual, 6)
	})
}

var j int32

func TestCustomWrap(t *testing.T) {
	Convey("CustomWrap", t, func() {
		wg := NewGoGroup("group")
		wg.Wrap("+3", func() {
			atomic.AddInt32(&j, 3)
		})
		wg.Wrap("+6", func() {
			atomic.AddInt32(&j, 6)
		})
		wg.Wrap("+9", func() {
			atomic.AddInt32(&j, 9)
		})
		wg.Wait()
		So(j, ShouldEqual, 18)
	})
}

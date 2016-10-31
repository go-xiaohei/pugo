package helper

import (
	"errors"
	"sync/atomic"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var i int32

func TestGlobalWrap(t *testing.T) {
	Convey("GlobalWrap", t, func() {
		GoWrap("+1", func() error {
			atomic.AddInt32(&i, 1)
			return nil
		})
		GoWrap("+2", func() error {
			atomic.AddInt32(&i, 2)
			return errors.New("2")
		})
		GoWrap("+3", func() error {
			atomic.AddInt32(&i, 3)
			return nil
		})
		GoWait()
		So(i, ShouldEqual, 6)
		So(GoWrapErrors(), ShouldHaveLength, 1)
	})
}

var j int32

func TestCustomWrap(t *testing.T) {
	Convey("CustomWrap", t, func() {
		wg := NewGoGroup("group")
		wg.Wrap("+3", func() error {
			atomic.AddInt32(&j, 3)
			return nil
		})
		wg.Wrap("+6", func() error {
			atomic.AddInt32(&j, 6)
			return nil
		})
		wg.Wrap("+9", func() error {
			atomic.AddInt32(&j, 9)
			return nil
		})
		wg.Wait()
		So(j, ShouldEqual, 18)
	})
}

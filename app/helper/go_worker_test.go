package helper

import (
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGoWorker(t *testing.T) {
	Convey("GoWorker", t, func() {
		w := NewGoWorker()
		w.Start()
		req := &GoWorkerRequest{
			Ctx: context.Background(),
			Action: func(ctx context.Context) (context.Context, error) {
				ctx = context.WithValue(ctx, "worker", "worker")
				return ctx, nil
			},
		}
		w.Send(req)
		w.Send(req)
		w.Send(req)
		w.Receive(func(res *GoWorkerResult) {
			Convey("GoWorkerResult", t, func() {
				So(res.Ctx.Value("worker").(string), ShouldEqual, "worker")
			})
		})
		w.WaitStop()
	})

	Convey("GoworkerChan", t, func() {
		w := NewGoWorker()
		w.Start()
		req := &GoWorkerRequest{
			Ctx: context.Background(),
			Action: func(ctx context.Context) (context.Context, error) {
				ctx = context.WithValue(ctx, "worker", "worker")
				return ctx, nil
			},
		}
		w.Send(req)
		w.Send(req)
		w.Send(req)
		go func() {
			for {
				res := <-w.Result()
				if res == nil {
					break
				}
				func(res *GoWorkerResult) {
					Convey("GoWorkerResult", t, func() {
						So(res.Ctx.Value("worker").(string), ShouldEqual, "worker")
					})
				}(res)
			}
		}()
		w.WaitStop()
	})

	Convey("GoWorkerErrorCase", t, func() {
		defer func() {
			if err := recover(); err != nil {
				So(fmt.Sprint(err), ShouldEqual, "nil GoWorker.Receive function")
			}
		}()
		w := NewGoWorker()
		w.Start()
		w.Receive(nil)
		w.WaitStop()
	})

	Convey("GoWorkerErrorCase2", t, func() {
		defer func() {
			if err := recover(); err != nil {
				So(fmt.Sprint(err), ShouldEqual, "GoWorker.Receive was called")
			}
		}()
		w := NewGoWorker()
		w.Start()
		w.Receive(func(res *GoWorkerResult) {
			Convey("GoWorkerResult", t, func() {
				So(res.Ctx.Value("worker").(string), ShouldEqual, "worker")
			})
		})
		w.Receive(func(res *GoWorkerResult) {
			Convey("GoWorkerResult", t, func() {
				So(res.Ctx.Value("worker").(string), ShouldEqual, "worker")
			})
		})
		w.WaitStop()
	})
}

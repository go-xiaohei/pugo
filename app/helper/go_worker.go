package helper

import (
	"context"
	"runtime"
	"sync/atomic"
)

// GoWorker is a goroutine pool to execute tasks
type GoWorker struct {
	channels   []chan *GoWorkerRequest
	resultChan chan *GoWorkerResult
	step       uint32
}

// GoWorkerRequest is goroutine worker task
type GoWorkerRequest struct {
	Ctx    context.Context
	Action func(ctx context.Context) (context.Context, error)
}

// GoWorkerResult is task result in goroutine worker
type GoWorkerResult struct {
	Ctx   context.Context
	Error error
}

// NewGoWorker create new goroutine worker with NumCPU goroutines
func NewGoWorker() *GoWorker {
	w := &GoWorker{
		channels:   make([]chan *GoWorkerRequest, runtime.NumCPU()),
		resultChan: make(chan *GoWorkerResult, runtime.NumCPU()),
		step:       0,
	}
	for i := range w.channels {
		w.channels[i] = make(chan *GoWorkerRequest)
	}
	return w
}

// Start start goroutine tasks listening
func (gw *GoWorker) Start() {
	for i := range gw.channels {
		go func(c chan *GoWorkerRequest) {
			for {
				req := <-c
				if req.Action == nil && req.Ctx == nil {
					return
				}
				//ctx = context.WithCancel()
				//ctx = context.WithDeadline()
				//ctx = context.WithTimeout()
				res := &GoWorkerResult{}
				res.Ctx, res.Error = req.Action(req.Ctx)
				gw.resultChan <- res
			}
		}(gw.channels[i])
	}
}

// Stop send stop signal to living goroutines
func (gw *GoWorker) Stop() {
	for i := range gw.channels {
		gw.channels[i] <- &GoWorkerRequest{
			Ctx:    nil,
			Action: nil,
		}
	}
}

// Send send task to goroutine in order, one by one
func (gw *GoWorker) Send(req *GoWorkerRequest) {
	step := atomic.LoadUint32(&gw.step)
	idx := int(step) % len(gw.channels)
	c := gw.channels[idx]
	c <- req
	atomic.AddUint32(&gw.step, 1)
}

// Recieve add result handler
func (gw *GoWorker) Recieve(fn func(res *GoWorkerResult)) {
	if fn == nil {
		panic("nil GoWorker.Recieve function")
	}
	go func() {
		for {
			res := <-gw.resultChan
			if res == nil {
				return
			}
			fn(res)
		}
	}()
}

// Result return result channel to handle task result in manual
func (gw *GoWorker) Result() chan *GoWorkerResult {
	return gw.resultChan
}

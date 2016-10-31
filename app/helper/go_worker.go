package helper

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
)

// GoWorker is a goroutine pool to execute tasks
type GoWorker struct {
	channels   []chan *GoWorkerRequest
	resultChan chan *GoWorkerResult
	step       uint32
	wg         sync.WaitGroup
	isReceived bool
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
		gw.wg.Add(1)
		go func(c chan *GoWorkerRequest) {
			for {
				req := <-c
				if req.Action == nil && req.Ctx == nil {
					break
				}
				//ctx = context.WithCancel()
				//ctx = context.WithDeadline()
				//ctx = context.WithTimeout()
				res := &GoWorkerResult{}
				res.Ctx, res.Error = req.Action(req.Ctx)
				gw.resultChan <- res
			}
			gw.wg.Done()
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

// WaitStop stop and wait all goroutine Done
func (gw *GoWorker) WaitStop() {
	gw.Stop()
	gw.wg.Wait()
	gw.resultChan <- nil // close result channel goroutine
}

// Send send task to goroutine in order, one by one
func (gw *GoWorker) Send(req *GoWorkerRequest) {
	step := atomic.LoadUint32(&gw.step)
	idx := int(step) % len(gw.channels)
	c := gw.channels[idx]
	c <- req
	atomic.AddUint32(&gw.step, 1)
}

// Receive add result handler
func (gw *GoWorker) Receive(fn func(res *GoWorkerResult)) {
	if fn == nil {
		panic("nil GoWorker.Receive function")
	}
	if gw.isReceived {
		panic("GoWorker.Receive was called")
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
	gw.isReceived = true
}

// Result return result channel to handle task result in manual
func (gw *GoWorker) Result() chan *GoWorkerResult {
	return gw.resultChan
}

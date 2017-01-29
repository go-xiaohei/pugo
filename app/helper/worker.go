package helper

import (
	"runtime"
	"sync"
)

type (
	// Worker is a worker pool for function
	Worker struct {
		funcs     []WorkerFunc
		chans     []chan WorkerFunc
		wg        sync.WaitGroup
		errors    []error
		errorLock sync.Mutex
	}
	// WorkerFunc is handler in Worker
	WorkerFunc func() error
)

// NewWorker creates new Worker with size.
// If size is 0, use runtime.NumCPU()
func NewWorker(size int) *Worker {
	w := &Worker{}
	if size == 0 {
		size = runtime.NumCPU()
	}
	w.chans = make([]chan WorkerFunc, size)
	for i := range w.chans {
		w.chans[i] = make(chan WorkerFunc)
	}
	w.startLoop()
	return w
}

func (w *Worker) startLoop() {
	w.wg.Add(len(w.chans))
	for _, ch := range w.chans {
		go func(ch chan WorkerFunc) {
			for {
				w.wg.Add(1)
				fn := <-ch
				if fn == nil {
					w.wg.Done()
					break
				}
				if err := fn(); err != nil {
					w.errorLock.Lock()
					w.errors = append(w.errors, err)
					w.errorLock.Unlock()
				}
				w.wg.Done()
			}
			w.wg.Done()
		}(ch)
	}
}

// AddFunc adds WorkerFunc
func (w *Worker) AddFunc(fn WorkerFunc) {
	w.funcs = append(w.funcs, fn)
}

// RunOnce runs WorkerFunc and waits all goroutine end
func (w *Worker) RunOnce() {
	size := len(w.chans)
	for i, fn := range w.funcs {
		idx := i % size
		w.chans[idx] <- fn
	}
	for _, ch := range w.chans {
		ch <- nil
	}
	w.wg.Wait()
}

// Errors returns errors when running WorkerFunc
func (w *Worker) Errors() []error {
	return w.errors
}

package builder

import (
	"time"

	"gopkg.in/inconshreveable/log15.v2"
)

type (
	// Builder is object of Builder handlers
	Builder struct {
		IsBuilding bool
		IsWatching bool
		Counter    int

		handlers []Handler
	}
	// Handler define a step in building process
	Handler func(ctx *Context)
)

var (
	b = &Builder{
		IsBuilding: false,
		IsWatching: false,
		handlers: []Handler{
			ReadSource,
			ReadTheme,
			AssembleSource,
			Compile,
			Copy,
		},
	}
	b2 = &Builder{
		IsBuilding: false,
		IsWatching: false,
		handlers: []Handler{
			ReadSource,
		},
	}
)

// Before add handler before building
func Before(fn Handler) {
	b.handlers = append([]Handler{fn}, b.handlers...)
}

// After add handler after building
func After(fn Handler) {
	b.handlers = append(b.handlers, fn)
}

// Build do a process with Context.
// the context should be prepared before building.
func Build(ctx *Context) {
	b.IsBuilding = true
	t := time.Now()
	log15.Info("Build|Start")
	for _, h := range b.handlers {
		if h(ctx); ctx.Err != nil {
			log15.Crit("Build|Fail|%s", ctx.Err.Error())
			break
		}
		log15.Debug("Build|Handler|%.3fms", time.Since(t).Seconds()*1e3)
	}
	log15.Info("Build|%d Pages", ctx.counter)
	b.IsBuilding = false
	b.Counter++
	if ctx.Err == nil {
		log15.Info("Build|Done|%d|%.1fms", Counter(), ctx.Duration()*1e3)
	}
}

// Counter return the times of building process ran
func Counter() int {
	return b.Counter
}

// Read do a process to read Source with Context.
// It does not build any thing, just read source data.
func Read(ctx *Context) {
	for _, h := range b2.handlers {
		if h(ctx); ctx.Err != nil {
			log15.Crit("Read|Fail|%s", ctx.Err.Error())
			break
		}
	}
}

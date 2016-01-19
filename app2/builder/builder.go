package builder

import "gopkg.in/inconshreveable/log15.v2"

type (
	// Builder is object of Builder handlers
	Builder struct {
		IsBuilding bool
		IsWatching bool

		handlers []Handler
	}
	// Handler define a step in building process
	Handler func(ctx *Context)
)

var (
	b = new()
)

func new() *Builder {
	return &Builder{
		IsBuilding: false,
		IsWatching: false,
		handlers:   []Handler{},
	}
}

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
	log15.Info("Build|Start")
	for _, h := range b.handlers {
		if h(ctx); ctx.Err != nil {
			log15.Crit("Build|Fail|%s", ctx.Err.Error)
		}
	}
	log15.Info("Build|Done|%s|%s", "a", "b")
}

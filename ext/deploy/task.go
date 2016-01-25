package deploy

import "github.com/go-xiaohei/pugo/app/builder"

var (
	manager *Manager
)

func init() {
	manager = NewManager(new(Git))
}

type (
	Task interface {
		Name() string
		Detect(*builder.Context) (Task, bool)
		Action(*builder.Context) error
	}
	Manager struct {
		tasks map[string]Task
	}
)

func NewManager(tasks ...Task) *Manager {
	m := &Manager{
		tasks: make(map[string]Task),
	}
	for _, t := range tasks {
		m.tasks[t.Name()] = t
	}
	return m
}

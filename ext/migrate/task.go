package migrate

import "github.com/go-xiaohei/pugo/app/builder"

var (
	manager *Manager
)

func init() {
	manager = NewManager()
}

type (
	Task interface {
		Name() string
		Detect(*builder.Context) (Task, error)
		Action(*builder.Context) error
	}
	// Manager manage tasks in global
	Manager struct {
		tasks map[string]Task
	}
)

// NewManager new manager with typed tasks
func NewManager(tasks ...Task) *Manager {
	m := &Manager{
		tasks: make(map[string]Task),
	}
	for _, t := range tasks {
		m.tasks[t.Name()] = t
	}
	return m
}

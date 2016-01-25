package deploy

import "github.com/go-xiaohei/pugo/app/builder"

var (
	manager *Manager
)

func init() {
	manager = NewManager(new(Git))
}

type (
	// Task define deploy task interface
	Task interface {
		Name() string
		Detect(*builder.Context) (Task, error)
		Action(*builder.Context) error
	}
	// Manager manage tasks in global
	Manager struct {
		tasks map[string]Task
		task  Task
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

// Set set using task
func (m *Manager) Set(task Task) {
	m.task = task
}

// Get get using task
func (m *Manager) Get() Task {
	return m.task
}

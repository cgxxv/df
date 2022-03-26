package dag

import (
	"context"
	"sync"
)

type Bus interface {
	sync.Locker
}

type Task interface {
	Process(ctx context.Context, bus Bus) error
}

type task struct {
	Task
	Name string

	waits    map[string]struct{}
	notifies map[string]struct{}
	waits2   map[string]struct{}
}

func NewTask(name string) *task {
	return &task{
		Name:     name,
		waits:    make(map[string]struct{}, 16),
		notifies: make(map[string]struct{}, 16),
		waits2:   make(map[string]struct{}, 16),
	}
}

func (t *task) Process(ctx context.Context, bus Bus) error {
	if t.Task != nil {
		bus.Lock()
		defer bus.Unlock()
		return t.Task.Process(ctx, bus)
	}
	return nil
}

func (t *task) Wait(t1 *task) {
	t.waits[t1.Name] = struct{}{}
	t1.notifies[t.Name] = struct{}{}
}

func (t *task) Reset() {
	for k, v := range t.waits {
		t.waits2[k] = v
	}
}

func (t *task) Waits() map[string]struct{} {
	waits := make(map[string]struct{}, len(t.waits))
	for k, v := range t.waits {
		waits[k] = v
	}

	return waits
}

func (t *task) String() string {
	if t.Name != "" {
		return t.Name
	}

	return "<nil>"
}

package dag

import (
	"context"
	"errors"
)

type Task interface {
	GetName() string
	Process(ctx context.Context, bus any) error
}

var (
	taskMap = make(map[string]Task, 1024)
)

func RegisterTask(task ...Task) error {
	for _, t := range task {
		_, ok := taskMap[t.GetName()]
		if ok {
			return errors.New("Duplicated task name")
		}
		taskMap[t.GetName()] = t
	}
	return nil
}

type task struct {
	Task

	waits    map[string]struct{}
	notifies map[string]struct{}
	waits2   map[string]struct{}
	err      error
}

func NewTask(ut Task) *task {
	return &task{
		Task:     ut,
		waits:    make(map[string]struct{}, 16),
		notifies: make(map[string]struct{}, 16),
		waits2:   make(map[string]struct{}, 16),
		err:      nil,
	}
}

func (t *task) Process(ctx context.Context, bus any) error {
	if t.Task != nil {
		return t.Task.Process(ctx, bus)
	}
	return nil
}

func (t *task) Wait(t1 *task) {
	t.waits[t1.GetName()] = struct{}{}
	t1.notifies[t.GetName()] = struct{}{}
}

func (t *task) Reset() {
	for k, v := range t.waits {
		t.waits2[k] = v
	}
	t.err = nil
}

func (t *task) Waits() map[string]struct{} {
	waits := make(map[string]struct{}, len(t.waits))
	for k, v := range t.waits {
		waits[k] = v
	}

	return waits
}

func (t *task) String() string {
	return t.GetName()
}

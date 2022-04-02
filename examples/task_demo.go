package main

import (
	"context"
	"sync/atomic"
	"time"
)

type MyBus struct {
	Data int32
}

type timeout struct{}

func (t *timeout) Timeout() time.Duration { return 1 * time.Millisecond }

type TaskDemo1 struct {
	timeout
}

func (t *TaskDemo1) GetName() string { return "1" }

func (t *TaskDemo1) Process(ctx context.Context, bus any) error {
	atomic.AddInt32(&bus.(*MyBus).Data, 1)

	return nil
}

type TaskDemo2 struct {
	timeout
}

func (t *TaskDemo2) GetName() string { return "2" }

func (t *TaskDemo2) Process(ctx context.Context, bus any) error {
	atomic.AddInt32(&bus.(*MyBus).Data, 1)

	return nil
}

type TaskDemo3 struct {
	timeout
}

func (t *TaskDemo3) GetName() string { return "3" }

func (t *TaskDemo3) Process(ctx context.Context, bus any) error {
	atomic.AddInt32(&bus.(*MyBus).Data, 1)

	return nil
}

type TaskDemo4 struct {
	timeout
}

func (t *TaskDemo4) GetName() string { return "4" }

func (t *TaskDemo4) Process(ctx context.Context, bus any) error {
	atomic.AddInt32(&bus.(*MyBus).Data, 1)

	return nil
}

type TaskDemo5 struct {
	timeout
}

func (t *TaskDemo5) GetName() string { return "5" }

func (t *TaskDemo5) Process(ctx context.Context, bus any) error {
	atomic.AddInt32(&bus.(*MyBus).Data, 1)

	return nil
}

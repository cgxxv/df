package dag

import (
	"context"
	"errors"
	"time"
)

type Task interface {
	GetName() string
	Timeout() time.Duration
	Process(ctx context.Context, bus any) error
}

var (
	taskMap = make(map[string]Task, 1024)
)

func RegisterTask(task ...Task) error {
	for _, t := range task {
		_, ok := taskMap[t.GetName()]
		if ok {
			return errors.New("Duplicated node name")
		}
		taskMap[t.GetName()] = t
	}
	return nil
}

package dag

import "context"

type node struct {
	uuid     string
	waits    map[string]struct{}
	notifies map[string]struct{}
	waits2   map[string]struct{}
	err      error
	done     chan struct{}

	Task
}

func NewNode(uuid string, ut Task) *node {
	return &node{
		uuid:     uuid,
		waits:    make(map[string]struct{}, 16),
		notifies: make(map[string]struct{}, 16),
		waits2:   make(map[string]struct{}, 16),
		err:      nil,
		done:     make(chan struct{}),

		Task: ut,
	}
}

func (t *node) Process(ctx context.Context, bus any) error {
	if t.Task != nil {
		return t.Task.Process(ctx, bus)
	}
	return nil
}

func (t *node) Wait(t1 *node) {
	t.waits[t1.GetName()] = struct{}{}
	t1.notifies[t.GetName()] = struct{}{}
}

func (t *node) Reset() {
	for k, v := range t.waits {
		t.waits2[k] = v
	}
	t.err = nil
}

func (t *node) Waits() map[string]struct{} {
	waits := make(map[string]struct{}, len(t.waits))
	for k, v := range t.waits {
		waits[k] = v
	}

	return waits
}

func (t *node) String() string {
	return t.GetName()
}

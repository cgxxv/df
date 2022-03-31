package dag

import (
	"container/list"
	"context"
	"errors"
	"os"
	"reflect"
	"sync"

	"github.com/modern-go/reflect2"
)

var (
	dagmap = sync.Map{}
)

type daging struct {
	dagpool  *sync.Pool
	nodemeta map[string]*NodeMeta
	endnodes []*NodeMeta
}

func InitScheduler(ctx context.Context, service string) error {
	if _, ok := dagmap.Load(service); ok {
		return errors.New("Duplicated key for daging")
	}

	nms := getSchedulerNodeMeta(ctx, service)

	if len(nms) == 0 {
		return errors.New("Empty node lists")
	}

	dg := daging{
		nodemeta: make(map[string]*NodeMeta, len(nms)),
		endnodes: make([]*NodeMeta, 0, len(nms)),
	}

	var start = make([]*NodeMeta, 0, len(nms))

	for _, nm := range nms {
		if nm.Application != os.Getenv("APPLICATION") {
			return errors.New("Unmatched application name from environment")
		}

		if len(start) > 1 {
			return errors.New("Start task must be one node")
		}

		if len(nm.Children) == 0 {
			dg.endnodes = append(dg.endnodes, nm)
		} else if len(nm.Parents) == 0 {
			start = append(start, nm)
		}

		dg.nodemeta[nm.Name] = nm
	}

	if len(start) != 1 {
		return errors.New("Start task must be one node")
	}

	dg.dagpool = &sync.Pool{
		New: func() interface{} {
			d, err := dg.instance()
			if err != nil {
				return nil
			}

			return d
		},
	}

	dagmap.Store(service, dg)

	return nil
}

func Schedule(ctx context.Context, service string, bus Bus) error {
	if rt := reflect2.TypeOf(bus); rt.Kind() != reflect.Ptr {
		return errors.New("Bus must be a pointer")
	}
	di, ok := dagmap.Load(service)
	if !ok {
		return errors.New("Not found daging instance for service: " + service)
	}
	dg, ok := di.(daging)
	if !ok {
		panic("unreachable")
	}
	dag := dg.get()

	for _, v := range dag.SourceVertices() {
		go dag.execTask(ctx, v.task, bus)
	}

	var finished int

	for {
		if finished == dag.Size() {
			break
		}

		t := <-dag.TaskChan()
		if t.err != nil {
			return t.err
		}
		finished++

		for n := range t.notifies {
			w, ok := dag.vertmap[n]
			if !ok {
				panic("unreachable")
			}
			delete(w.task.waits2, t.GetName())
			if len(w.task.waits2) == 0 {
				go dag.execTask(ctx, w.task, bus)
			}
		}
	}

	_ = dg.put(dag)

	return nil
}

func (dg *daging) instance() (*DAG, error) {
	if len(dg.endnodes) <= 0 {
		return nil, errors.New("Empty end nodemeta")
	}
	if len(dg.nodemeta) <= 0 {
		return nil, errors.New("Empty nodemeta info")
	}

	d := NewDAG()
	q := list.New()
	for _, endnode := range dg.endnodes {
		q.PushBack(endnode)
	}

	for q.Len() > 0 {
		e := q.Front()
		nm, ok := e.Value.(*NodeMeta)
		if !ok {
			return nil, errors.New("Node assertion error")
		}

		var vertex *Vertex
		if vt, ok := d.vertmap[nm.Name]; !ok {
			ut, ok := taskMap[nm.Name]
			if !ok {
				panic("unreachable")
			}
			vertex = NewVertex(nm.Name, NewTask(ut))
			d.AddVertex(vertex)
		} else {
			vertex = vt
		}

		for _, parent := range nm.Parents {
			q.PushBack(dg.nodemeta[parent])

			var vt *Vertex
			if _vt, ok := d.vertmap[parent]; !ok {
				ut, ok := taskMap[parent]
				if !ok {
					panic("unreachable")
				}
				vt = NewVertex(parent, NewTask(ut))
				d.AddVertex(vt)
			} else {
				vt = _vt
			}
			d.AddEdge(vt, vertex)
		}
		q.Remove(e)
	}

	return d, nil
}

func (dg *daging) get() *DAG {
	d, ok := dg.dagpool.Get().(*DAG)
	if !ok {
		panic("unreachable")
	}
	d.reset()

	return d
}

func (dg *daging) put(d *DAG) error {
	defer dg.dagpool.Put(d)
	return d.reset()
}

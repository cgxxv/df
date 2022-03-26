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
	endnode  *NodeMeta
}

func InitNodeMeta(service string, nms []*NodeMeta) error {
	if _, ok := dagmap.Load(service); ok {
		return errors.New("Duplicated key for daging")
	}

	if len(nms) == 0 {
		return errors.New("Empty node lists")
	}

	dg := daging{
		nodemeta: make(map[string]*NodeMeta, len(nms)),
	}

	var start = make([]*NodeMeta, 0, len(nms))
	var end = make([]*NodeMeta, 0, len(nms))

	for _, nm := range nms {
		if nm.Application != os.Getenv("APPLICATION") {
			return errors.New("Unmatched application name from environment")
		}

		if len(start) > 1 {
			return errors.New("Start task must be one node")
		}
		if len(end) > 1 {
			return errors.New("End task must be one node")
		}

		if len(nm.Children) == 0 {
			end = append(end, nm)
		} else if len(nm.Parent) == 0 {
			start = append(start, nm)
		}

		dg.nodemeta[nm.Name] = nm
	}

	if len(start) != 1 {
		return errors.New("Start task must be one node")
	}
	if len(end) != 1 {
		return errors.New("End task must be one node")
	}

	dg.endnode = end[0]
	dg.dagpool = &sync.Pool{
		New: func() interface{} {
			d, err := dg.NewDag()
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
		go runTask(ctx, dag, v.task, bus)
	}

	var finished int

	for {
		if finished == dag.Size() {
			break
		}

		t := <-dag.TaskChan()
		finished++

		for n := range t.notifies {
			w, ok := dag.vertmap[n]
			if !ok {
				panic("unreachable")
			}
			delete(w.task.waits2, t.Name)
			if len(w.task.waits2) == 0 {
				go runTask(ctx, dag, w.task, bus)
			}
		}
	}

	_ = dg.put(dag)

	return nil
}

func runTask(ctx context.Context, dag *DAG, t *task, bus Bus) {
	if err := t.Process(ctx, bus); err != nil {
		println(err)
	}
	dag.TaskChan() <- t
}

func (dg *daging) NewDag() (*DAG, error) {
	if dg.endnode == nil {
		return nil, errors.New("Nil end nodemeta")
	}
	if len(dg.nodemeta) <= 0 {
		return nil, errors.New("Empty nodemeta info")
	}

	d := NewDAG()
	q := list.New()
	q.PushBack(dg.endnode)

	for q.Len() > 0 {
		e := q.Front()
		nm, ok := e.Value.(*NodeMeta)
		if !ok {
			return nil, errors.New("Node assertion error")
		}

		var vertex *Vertex
		if vt, ok := d.vertmap[nm.Name]; !ok {
			t := NewTask(nm.Name)
			t.Task = nodeMap[nm.Name]
			vertex = NewVertex(nm.Name, t)
			d.AddVertex(vertex)
		} else {
			vertex = vt
		}

		for _, parent := range nm.Parent {
			q.PushBack(dg.nodemeta[parent])

			var vt *Vertex
			if _vt, ok := d.vertmap[parent]; !ok {
				pt := NewTask(parent)
				pt.Task = nodeMap[parent]
				vt = NewVertex(parent, pt)
				d.AddVertex(vt)
			} else {
				vt = _vt
			}
			addEdge(d, vt, vertex)
		}
		q.Remove(e)
	}

	return d, nil
}

func addEdge(dag *DAG, tailVertex, headVertex *Vertex) {
	headVertex.task.Wait(tailVertex.task)
	dag.AddEdge(tailVertex, headVertex)
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

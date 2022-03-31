package dag

import (
	"context"
	"errors"
	"fmt"

	"github.com/goombaio/orderedmap"
)

type DAG struct {
	vertices orderedmap.OrderedMap
	vertmap  map[string]*Vertex
	taskChan chan *task
}

func NewDAG() *DAG {
	d := &DAG{
		vertices: *orderedmap.NewOrderedMap(),
		vertmap:  make(map[string]*Vertex, 256),
		taskChan: make(chan *task),
	}

	return d
}

func (d *DAG) TaskChan() chan *task {
	return d.taskChan
}

func (d *DAG) AddVertex(v *Vertex) error {
	d.vertices.Put(v.ID, v)
	d.vertmap[v.task.GetName()] = v

	return nil
}

func (d *DAG) DeleteVertex(vertex *Vertex) error {
	existsVertex := false

	for _, v := range d.vertices.Values() {
		if v == vertex {
			existsVertex = true
		}
	}
	if !existsVertex {
		return fmt.Errorf("Vertex with ID %v not found", vertex.ID)
	}

	d.vertices.Remove(vertex.ID)
	delete(d.vertmap, vertex.task.GetName())

	return nil
}

func (d *DAG) AddEdge(tailVertex *Vertex, headVertex *Vertex) error {
	tailExists := false
	headExists := false

	for _, vertex := range d.vertices.Values() {
		if vertex == tailVertex {
			tailExists = true
		}
		if vertex == headVertex {
			headExists = true
		}
	}
	if !tailExists {
		return fmt.Errorf("Vertex with ID %v not found", tailVertex.ID)
	}
	if !headExists {
		return fmt.Errorf("Vertex with ID %v not found", headVertex.ID)
	}

	for _, childVertex := range tailVertex.Children.Values() {
		if childVertex == headVertex {
			return fmt.Errorf("Edge (%v,%v) already exists", tailVertex.ID, headVertex.ID)
		}
	}

	tailVertex.Children.Add(headVertex)
	headVertex.Parents.Add(tailVertex)

	headVertex.task.Wait(tailVertex.task)

	return nil
}

func (d *DAG) DeleteEdge(tailVertex *Vertex, headVertex *Vertex) error {
	for _, childVertex := range tailVertex.Children.Values() {
		if childVertex == headVertex {
			tailVertex.Children.Remove(childVertex)
		}
	}

	return nil
}

func (d *DAG) GetVertex(id interface{}) (*Vertex, error) {
	var vertex *Vertex

	v, found := d.vertices.Get(id)
	if !found {
		return vertex, fmt.Errorf("vertex %s not found in the graph", id)
	}

	vertex = v.(*Vertex)

	return vertex, nil
}

func (d *DAG) Order() int {
	numVertices := d.vertices.Size()

	return numVertices
}

func (d *DAG) Size() int {
	numEdges := 0
	for _, vertex := range d.vertices.Values() {
		numEdges = numEdges + vertex.(*Vertex).Children.Size()
	}

	return numEdges
}

func (d *DAG) SinkVertices() []*Vertex {
	var sinkVertices []*Vertex

	for _, vertex := range d.vertices.Values() {
		if vertex.(*Vertex).Children.Size() == 0 {
			sinkVertices = append(sinkVertices, vertex.(*Vertex))
		}
	}

	return sinkVertices
}

func (d *DAG) SourceVertices() []*Vertex {
	var sourceVertices []*Vertex

	for _, vertex := range d.vertices.Values() {
		if vertex.(*Vertex).Parents.Size() == 0 {
			sourceVertices = append(sourceVertices, vertex.(*Vertex))
		}
	}

	return sourceVertices
}

func (d *DAG) Successors(vertex *Vertex) ([]*Vertex, error) {
	var successors []*Vertex

	_, found := d.GetVertex(vertex.ID)
	if found != nil {
		return successors, fmt.Errorf("vertex %s not found in the graph", vertex.ID)
	}

	for _, v := range vertex.Children.Values() {
		successors = append(successors, v.(*Vertex))
	}

	return successors, nil
}

func (d *DAG) Predecessors(vertex *Vertex) ([]*Vertex, error) {
	var predecessors []*Vertex

	_, found := d.GetVertex(vertex.ID)
	if found != nil {
		return predecessors, fmt.Errorf("vertex %s not found in the graph", vertex.ID)
	}

	for _, v := range vertex.Parents.Values() {
		predecessors = append(predecessors, v.(*Vertex))
	}

	return predecessors, nil
}

func (d *DAG) String() string {
	result := fmt.Sprintf("DAG Vertices: %d - Edges: %d\n", d.Order(), d.Size())
	result += fmt.Sprintf("Vertices:\n")
	for _, vertex := range d.vertices.Values() {
		vertex = vertex.(*Vertex)
		result += fmt.Sprintf("%s", vertex)
	}

	return result
}

func (d *DAG) execTask(ctx context.Context, t *task, bus Bus) {
	defer func() {
		if a := recover(); a != nil {
			switch err := a.(type) {
			case nil:
			case error:
				t.err = err
			default:
				t.err = fmt.Errorf("%v", err)
			}
		}
		d.TaskChan() <- t
	}()
	if err := t.Process(ctx, bus); err != nil {
		t.err = err
	}
}

func (d *DAG) reset() error {
	for _, v := range d.vertices.Values() {
		vertex, ok := v.(*Vertex)
		if !ok {
			return errors.New("not a Vertex")
		}
		vertex.task.Reset()
	}

	return nil
}

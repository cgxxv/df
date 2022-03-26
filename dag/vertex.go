package dag

import (
	"fmt"

	"github.com/goombaio/orderedset"
)

type Vertex struct {
	ID       string
	Parents  *orderedset.OrderedSet
	Children *orderedset.OrderedSet
	task     *task
}

func NewVertex(id string, task *task) *Vertex {
	v := &Vertex{
		ID:       id,
		Parents:  orderedset.NewOrderedSet(),
		Children: orderedset.NewOrderedSet(),
		task:     task,
	}

	return v
}

func (v *Vertex) Degree() int {
	return v.Parents.Size() + v.Children.Size()
}

func (v *Vertex) InDegree() int {
	return v.Parents.Size()
}

func (v *Vertex) OutDegree() int {
	return v.Children.Size()
}

func (v *Vertex) Task() *task {
	return v.task
}

func (v *Vertex) String() string {
	result := fmt.Sprintf("ID: %s - Parents: %d - Children: %d - Value: %v\n",
		v.ID, v.Parents.Size(), v.Children.Size(), fmt.Sprintf("name:%#v, waits:%#v, waits2:%#v", v.task.Name, v.task.waits, v.task.waits2))

	return result
}

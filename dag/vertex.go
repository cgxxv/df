package dag

import (
	"fmt"

	"github.com/goombaio/orderedset"
)

type Vertex struct {
	ID       string
	Parents  *orderedset.OrderedSet
	Children *orderedset.OrderedSet
	node     *node
}

func NewVertex(id string, node *node) *Vertex {
	v := &Vertex{
		ID:       id,
		Parents:  orderedset.NewOrderedSet(),
		Children: orderedset.NewOrderedSet(),
		node:     node,
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

func (v *Vertex) Task() Task {
	return v.node.Task
}

func (v *Vertex) String() string {
	result := fmt.Sprintf("ID: %s - Parents: %d - Children: %d - Value:%v\n",
		v.ID, v.Parents.Size(), v.Children.Size(), v.node.GetName())

	return result
}

package dag_test

import (
	"fmt"

	"github.com/cgxxv/df/dag"
)

func ExampleDAG_vertices() {
	dag1 := dag.NewDAG()

	vertex1 := dag.NewVertex("1", dag.NewTask(newFakeTask[T]("")))
	vertex2 := dag.NewVertex("2", dag.NewTask(newFakeTask[T]("")))
	vertex3 := dag.NewVertex("3", dag.NewTask(newFakeTask[T]("")))
	vertex4 := dag.NewVertex("4", dag.NewTask(newFakeTask[T]("")))

	err := dag1.AddVertex(vertex1)
	if err != nil {
		fmt.Printf("Can't add vertex to DAG: %s", err)
		panic(err)
	}
	err = dag1.AddVertex(vertex2)
	if err != nil {
		fmt.Printf("Can't add vertex to DAG: %s", err)
		panic(err)
	}
	err = dag1.AddVertex(vertex3)
	if err != nil {
		fmt.Printf("Can't add vertex to DAG: %s", err)
		panic(err)
	}
	err = dag1.AddVertex(vertex4)
	if err != nil {
		fmt.Printf("Can't add vertex to DAG: %s", err)
		panic(err)
	}

	fmt.Println(dag1.String())
	// Output:
	// DAG Vertices: 4 - Edges: 0
	// Vertices:
	// ID: 1 - Parents: 0 - Children: 0 - Value:
	// ID: 2 - Parents: 0 - Children: 0 - Value:
	// ID: 3 - Parents: 0 - Children: 0 - Value:
	// ID: 4 - Parents: 0 - Children: 0 - Value:
}

func ExampleDAG_edges() {
	dag1 := dag.NewDAG()

	vertex1 := dag.NewVertex("1", dag.NewTask(newFakeTask[T]("")))
	vertex2 := dag.NewVertex("2", dag.NewTask(newFakeTask[T]("")))
	vertex3 := dag.NewVertex("3", dag.NewTask(newFakeTask[T]("")))
	vertex4 := dag.NewVertex("4", dag.NewTask(newFakeTask[T]("")))

	err := dag1.AddVertex(vertex1)
	if err != nil {
		fmt.Printf("Can't add vertex to DAG: %s", err)
		panic(err)
	}
	err = dag1.AddVertex(vertex2)
	if err != nil {
		fmt.Printf("Can't add vertex to DAG: %s", err)
		panic(err)
	}
	err = dag1.AddVertex(vertex3)
	if err != nil {
		fmt.Printf("Can't add vertex to DAG: %s", err)
		panic(err)
	}
	err = dag1.AddVertex(vertex4)
	if err != nil {
		fmt.Printf("Can't add vertex to DAG: %s", err)
		panic(err)
	}

	// Edges

	err = dag1.AddEdge(vertex1, vertex2)
	if err != nil {
		fmt.Printf("Can't add edge to DAG: %s", err)
		panic(err)
	}

	err = dag1.AddEdge(vertex2, vertex3)
	if err != nil {
		fmt.Printf("Can't add edge to DAG: %s", err)
		panic(err)
	}

	err = dag1.AddEdge(vertex3, vertex4)
	if err != nil {
		fmt.Printf("Can't add edge to DAG: %s", err)
		panic(err)
	}

	fmt.Println(dag1.String())
	// Output:
	// DAG Vertices: 4 - Edges: 3
	// Vertices:
	// ID: 1 - Parents: 0 - Children: 1 - Value:
	// ID: 2 - Parents: 1 - Children: 1 - Value:
	// ID: 3 - Parents: 1 - Children: 1 - Value:
	// ID: 4 - Parents: 1 - Children: 0 - Value:
}

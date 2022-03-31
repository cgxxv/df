package dag_test

import (
	"testing"

	"github.com/cgxxv/df/dag"
)

func TestVertex(t *testing.T) {
	v := dag.NewVertex("1", dag.NewTask(nil))

	if v.ID == "" {
		t.Fatalf("Vertex ID expected to be not empty string.\n")
	}
	if v.Task() != nil {
		t.Fatalf("Vertex Value expected to be nil.\n")
	}
}

func TestVertex_Parents(t *testing.T) {
	v := dag.NewVertex("1", dag.NewTask(nil))

	numParents := v.Parents.Size()
	if numParents != 0 {
		t.Fatalf("Vertex Parents expected to be 0 but got %d", v.Parents.Size())
	}
}

func TestVertex_Children(t *testing.T) {
	v := dag.NewVertex("1", dag.NewTask(nil))

	numParents := v.Children.Size()
	if numParents != 0 {
		t.Fatalf("Vertex Children expected to be 0 but got %d", v.Children.Size())
	}
}

func TestVertex_Degree(t *testing.T) {
	dag1 := dag.NewDAG()

	vertex1 := dag.NewVertex("1", dag.NewTask(nil))
	vertex2 := dag.NewVertex("2", dag.NewTask(nil))
	vertex3 := dag.NewVertex("3", dag.NewTask(nil))

	degree := vertex1.Degree()
	if degree != 0 {
		t.Fatalf("Vertex1 Degree expected to be 0 but got %d", vertex1.Degree())
	}

	err := dag1.AddVertex(vertex1)
	if err != nil {
		t.Fatalf("Can't add vertex to DAG: %s", err)
	}
	err = dag1.AddVertex(vertex2)
	if err != nil {
		t.Fatalf("Can't add vertex to DAG: %s", err)
	}
	err = dag1.AddVertex(vertex3)
	if err != nil {
		t.Fatalf("Can't add vertex to DAG: %s", err)
	}

	err = dag1.AddEdge(vertex1, vertex2)
	if err != nil {
		t.Fatalf("Can't add edge to DAG: %s", err)
	}

	err = dag1.AddEdge(vertex2, vertex3)
	if err != nil {
		t.Fatalf("Can't add edge to DAG: %s", err)
	}

	degree = vertex1.Degree()
	if degree != 1 {
		t.Fatalf("Vertex1 Degree expected to be 1 but got %d", vertex1.Degree())
	}

	degree = vertex2.Degree()
	if degree != 2 {
		t.Fatalf("Vertex2 Degree expected to be 2 but got %d", vertex2.Degree())
	}
}

func TestVertex_InDegree(t *testing.T) {
	dag1 := dag.NewDAG()

	vertex1 := dag.NewVertex("1", dag.NewTask(nil))
	vertex2 := dag.NewVertex("2", dag.NewTask(nil))
	vertex3 := dag.NewVertex("3", dag.NewTask(nil))

	inDegree := vertex1.InDegree()
	if inDegree != 0 {
		t.Fatalf("Vertex1 InDegree expected to be 0 but got %d", vertex1.InDegree())
	}

	err := dag1.AddVertex(vertex1)
	if err != nil {
		t.Fatalf("Can't add vertex to DAG: %s", err)
	}
	err = dag1.AddVertex(vertex2)
	if err != nil {
		t.Fatalf("Can't add vertex to DAG: %s", err)
	}
	err = dag1.AddVertex(vertex3)
	if err != nil {
		t.Fatalf("Can't add vertex to DAG: %s", err)
	}

	err = dag1.AddEdge(vertex1, vertex2)
	if err != nil {
		t.Fatalf("Can't add edge to DAG: %s", err)
	}

	err = dag1.AddEdge(vertex2, vertex3)
	if err != nil {
		t.Fatalf("Can't add edge to DAG: %s", err)
	}

	inDegree = vertex1.InDegree()
	if inDegree != 0 {
		t.Fatalf("Vertex1 InDegree expected to be 0 but got %d", vertex1.InDegree())
	}

	inDegree = vertex2.InDegree()
	if inDegree != 1 {
		t.Fatalf("Vertex2 InDegree expected to be 1 but got %d", vertex2.InDegree())
	}
}

func TestVertex_OutDegree(t *testing.T) {
	dag1 := dag.NewDAG()

	vertex1 := dag.NewVertex("1", dag.NewTask(nil))
	vertex2 := dag.NewVertex("2", dag.NewTask(nil))
	vertex3 := dag.NewVertex("3", dag.NewTask(nil))

	outDegree := vertex1.OutDegree()
	if outDegree != 0 {
		t.Fatalf("Vertex1 OutDegree expected to be 0 but got %d", vertex1.OutDegree())
	}

	err := dag1.AddVertex(vertex1)
	if err != nil {
		t.Fatalf("Can't add vertex to DAG: %s", err)
	}
	err = dag1.AddVertex(vertex2)
	if err != nil {
		t.Fatalf("Can't add vertex to DAG: %s", err)
	}
	err = dag1.AddVertex(vertex3)
	if err != nil {
		t.Fatalf("Can't add vertex to DAG: %s", err)
	}

	err = dag1.AddEdge(vertex1, vertex2)
	if err != nil {
		t.Fatalf("Can't add edge to DAG: %s", err)
	}

	err = dag1.AddEdge(vertex2, vertex3)
	if err != nil {
		t.Fatalf("Can't add edge to DAG: %s", err)
	}

	outDegree = vertex1.OutDegree()
	if outDegree != 1 {
		t.Fatalf("Vertex1 OutDegree expected to be 1 but got %d", vertex1.OutDegree())
	}

	outDegree = vertex2.OutDegree()
	if outDegree != 1 {
		t.Fatalf("Vertex2 OutDegree expected to be 1 but got %d", vertex2.OutDegree())
	}

	outDegree = vertex3.OutDegree()
	if outDegree != 0 {
		t.Fatalf("Vertex2 OutDegree expected to be 0 but got %d", vertex3.OutDegree())
	}
}

func TestVertex_String(t *testing.T) {
	v := dag.NewVertex("1", dag.NewTask(nil))
	vstr := v.String()

	expected := "ID: 1 - Parents: 0 - Children: 0 - Value: <nil>\n"
	if vstr != expected {
		t.Fatalf("Vertex stringer expected to be %q but got %q\n", expected, vstr)
	}
}

func TestVertex_String_WithStringValue(t *testing.T) {
	v := dag.NewVertex("1", dag.NewTask(newFakeTask("one")))
	vstr := v.String()

	expected := "ID: 1 - Parents: 0 - Children: 0 - Value: one\n"
	if vstr != expected {
		t.Fatalf("Vertex stringer expected to be %q but got %q\n", expected, vstr)
	}
}

func TestVertex_WithStringValue(t *testing.T) {
	v := dag.NewVertex("1", dag.NewTask(newFakeTask("one")))

	if v.ID == "" {
		t.Fatalf("Vertex ID expected to be not empty string.\n")
	}
	if v.Task().GetName() != "one" {
		t.Fatalf("Vertex Value expected to be one.\n")
	}
}

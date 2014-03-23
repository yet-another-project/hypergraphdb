package iterator_test
import (
    "fmt"
    "testing"
    "github.com/yet-another-project/hypergraphdb/element"
    "github.com/yet-another-project/hypergraphdb/iterator"
)

func TestRecursiveDFSIteratorRecursivePost(t *testing.T) {
    g := element.NewGraph("g")
    a := g.NewSubGraph("a")
    b := a.NewMutualNeighbour("b")
    c := a.NewMutualNeighbour("c")
    d := a.NewMutualNeighbour("d")

    it := iterator.NewRecursiveDFS(a)
    testData := []*element.Node{b, c, d, a, nil}

    stream := it.Stream()
    for i := range testData {
        node := <-stream
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestRecursiveDFSIteratorSequentialDeep(t *testing.T) {
    g := element.NewGraph("g")
    a := g.NewSubGraph("a")
    b := a.NewNeighbour("b")
    c := b.NewNeighbour("c")
    d := c.NewNeighbour("d")

    it := iterator.NewRecursiveDFS(a)
    testData := []*element.Node{d, c, b, a, nil}

    stream := it.Stream()
    for i := range testData {
        node := <-stream
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestRecursiveDFSIteratorSequentialWithCycleGraph(t *testing.T) {
    g := element.NewGraph("g")
    a := g.NewSubGraph("a")
    b := g.NewSubGraph("b")
    c := g.NewSubGraph("c")

    a.ConnectNeighbour(b)
    a.ConnectNeighbour(c)
    b.ConnectNeighbour(c)
    b.ConnectNeighbour(a)
    c.ConnectNeighbour(a)
    c.ConnectNeighbour(b)

    it := iterator.NewRecursiveDFS(a)

    testData := []*element.Node{c, b, a, nil}

    stream := it.Stream()
    for i := range testData {
        node := <-stream
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestRecursiveDFSIteratorSimpleCircular(t *testing.T) {
    g := element.NewGraph("g")
    a := g.NewSubGraph("a")
    b := a.NewMutualNeighbour("b")
    it := iterator.NewRecursiveDFS(a)

    testData := []*element.Node{b, a, nil}

    stream := it.Stream()
    for i := range testData {
        node := <-stream
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
    it = iterator.NewRecursiveDFS(b)
    testData = []*element.Node{a, b, nil}

    stream = it.Stream()
    for i := range testData {
        node := <-stream
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestRecursiveDFSIteratorSingleNode(t *testing.T) {
    g := element.NewGraph("g")
    a := g.NewSubGraph("a")
    it := iterator.NewRecursiveDFS(a)
    testData := []*element.Node{a, nil}

    stream := it.Stream()
    for i := range testData {
        node := <-stream
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func BenchmarkFullyConnectedRecursiveDFS(b *testing.B) {
    size := 1900
    fmt.Println(size)
    _, graph := createFullyConnectedGraph(size)
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        it := iterator.NewRecursiveDFS(graph)
        stream := it.Stream()
        for ; true; <-stream {
           // fmt.Println(node)
        }
    }
}


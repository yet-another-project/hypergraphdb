package iterator_test
import (
    "fmt"
    "testing"
    "github.com/yet-another-project/hypergraphdb/element"
    "github.com/yet-another-project/hypergraphdb/iterator"
)

func TestLinearDFSIteratorSequentialPost(t *testing.T) {
    g := element.NewGraph("g")
    a := g.NewSubGraph("a")
    b := a.NewMutualNeighbour("b")
    c := a.NewMutualNeighbour("c")
    d := a.NewMutualNeighbour("d")

    it := iterator.NewLinearDFS(a)
    testData := []*element.Node{b, c, d, a, nil}

    for i := range testData {
        node := it.Next()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestLinearDFSIteratorSequentialDeep(t *testing.T) {
    g := element.NewGraph("g")
    a := g.NewSubGraph("a")
    b := a.NewNeighbour("b")
    c := b.NewNeighbour("c")
    d := c.NewNeighbour("d")

    it := iterator.NewLinearDFS(a)
    testData := []*element.Node{d, c, b, a, nil}

    for i := range testData {
        node := it.Next()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestLinearDFSIteratorSequentialWithCycleGraph(t *testing.T) {
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

    it := iterator.NewLinearDFS(a)

    testData := []*element.Node{c, b, a, nil}

    for i := range testData {
        node := it.Next()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestLinearDFSIteratorSimpleCircular(t *testing.T) {
    g := element.NewGraph("g")
    a := g.NewSubGraph("a")
    b := a.NewMutualNeighbour("b")
    it := iterator.NewLinearDFS(a)

    testData := []*element.Node{b, a, nil}

    for i := range testData {
        node := it.Next()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
    it = iterator.NewLinearDFS(b)
    testData = []*element.Node{a, b, nil}

    for i := range testData {
        node := it.Next()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestLinearDFSIteratorSingleNode(t *testing.T) {
    g := element.NewGraph("g")
    a := g.NewSubGraph("a")
    it := iterator.NewLinearDFS(a)
    testData := []*element.Node{a, nil}

    for i := range testData {
        node := it.Next()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func BenchmarkFullyConnectedLinearDFS(b *testing.B) {
    size := 1900
    fmt.Println(size)
    _, graph := createFullyConnectedGraph(size)
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        it := iterator.NewLinearDFS(graph)
        for node := it.Next(); node != nil; node = it.Next() {
           // fmt.Println(node)
        }
    }
}

func TestLinearDFSIteratorChannelSimple(t *testing.T) {
    g := element.NewGraph("g")
    a := g.NewSubGraph("a")
    b := a.NewMutualNeighbour("b")
    c := a.NewMutualNeighbour("c")
    d := a.NewMutualNeighbour("d")

    it := iterator.NewLinearDFS(a)
    go it.Run()
    testData := []*element.Node{b, c, d, a, nil}

    for i := range testData {
        node := <-it.Stream()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestLinearDFSIteratorChannelClose(t *testing.T) {
    g := element.NewGraph("g")
    a := g.NewSubGraph("a")
    b := a.NewMutualNeighbour("b")
    a.NewMutualNeighbour("c")
    a.NewMutualNeighbour("d")

    it := iterator.NewLinearDFS(a)
    go it.Run()
    node := <-it.Stream()
    if b != node {
        t.Error("DFS expected to deliver b, but instead", node)
    }

    it.Close()
}


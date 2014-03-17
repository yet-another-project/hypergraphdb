package graphdb
import (
    "testing"
    "fmt"
)

func TestNewGraph(t *testing.T) {
    g := NewGraph("test")
    if "test" != g.String() {
        t.Error("wrong value")
    }
}

func TestNewSubGraph(t *testing.T) {
    g := NewGraph("g")
    h := g.NewSubGraph("h")
    if "g[h]" != g.String() {
        t.Error("wrong value")
    }
    if h != g.subnodes[0] {
        t.Error("wrong subnode")
    }
}

func TestNewNeighbour(t *testing.T) {
    g := NewGraph("g")
    h := g.NewNeighbour("h")
    if nil != h {
        t.Error("expected nil")
    }
    m := g.NewSubGraph("m")
    n := m.NewNeighbour("n")
    if "g[m, n]" != g.String() {
        t.Error("actual " + g.String())
    }
    if "m{n}" != m.String() {
        t.Error("actual " + m.String())
    }
    if n != m.neighbours[0] {
        t.Error("wrong value")
    }
    if len(n.neighbours) != 0 {
        t.Error("n should not have neighbours")
    }
}

func TestNewMutualNeighbour(t *testing.T) {
    g := NewGraph("g")
    h := g.NewNeighbour("h")
    if nil != h {
        t.Error("expected nil")
    }
    m := g.NewSubGraph("m")
    n := m.NewMutualNeighbour("n")
    if "g[m, n]" != g.String() {
        t.Error("actual " + g.String())
    }
    if "m{n}" != m.String() {
        t.Error("actual " + m.String())
    }
    if "n{m}" != n.String() {
        t.Error("actual " + n.String())
    }
    if n != m.neighbours[0] {
        t.Error("wrong value")
    }
    if len(n.neighbours) != 1 {
        t.Error("n should have neighbours")
    }
    if n.neighbours[0] != m {
        t.Error("m should also be n's neighbour")
    }

    if nil != g.NewMutualNeighbour("z") {
        t.Error("expected nil")
    }
}

func TestConnectNeighbour(t *testing.T) {
    g := NewGraph("g")
    x := g.NewSubGraph("x")
    y := x.NewMutualNeighbour("y")
    z := y.NewNeighbour("z")
    if !z.ConnectNeighbour(x) {
        t.Error("expected success")
    }
    if !x.ConnectNeighbour(z) {
        t.Error("expected success")
    }
    if "z{x}" != z.String() {
        t.Error("expected " + z.String())
    }
    if x.ConnectNeighbour(y) {
        t.Error("expected failure")
    }
}

func TestUpwardParents(t *testing.T) {
    a := NewGraph("a")
    b := a.NewSubGraph("b")
    c := b.NewSubGraph("c")
    d := c.NewSubGraph("d")
    expected := NodeSet(nil)
    expected = append(expected, c)
    expected = append(expected, b)
    expected = append(expected, a)
    actual := d.UpwardParents()
    if len(actual) != len(expected) {
        t.Error("parents len differ, actual ", len(actual), " expected ", len(expected))
    }
    for i := range expected {
        if expected[i] != actual[i] {
            t.Error("parents differ at index", i, "expected", expected[i], "actual", actual[i])
        }
    }
}

func TestCommonAncestor(t *testing.T) {
    a := NewGraph("a")
    b := a.NewSubGraph("b")
    c := b.NewSubGraph("c")
    d := c.NewSubGraph("d")
    e := b.NewSubGraph("e")

    testData := [][]*Node{
        {b, e, d},
        {b, e, c},
        {a, e, b},
        {a, b, b},
        {nil, a, a},
    }
    for i := range testData {
        if testData[i][0] != testData[i][1].CommonAncestor(testData[i][2]) {
            t.Error(testData[i][0], "expected to be the common ancestor of", testData[i][1], "and", testData[i][2])
        }
        if testData[i][0] != testData[i][2].CommonAncestor(testData[i][1]) {
            t.Error(testData[i][0], "expected to be the common ancestor of", testData[i][2], "and", testData[i][1])
        }
    }
}

func TestDFSIteratorSequentialPost(t *testing.T) {
    g := NewGraph("g")
    a := g.NewSubGraph("a")
    b := a.NewMutualNeighbour("b")
    c := a.NewMutualNeighbour("c")
    d := a.NewMutualNeighbour("d")

    it := a.NewDFSIterator()
    testData := []*Node{b, c, d, a, nil}

    for i := range testData {
        node := it.Next()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestDFSIteratorSequentialDeep(t *testing.T) {
    g := NewGraph("g")
    a := g.NewSubGraph("a")
    b := a.NewNeighbour("b")
    c := b.NewNeighbour("c")
    d := c.NewNeighbour("d")

    it := a.NewDFSIterator()
    testData := []*Node{d, c, b, a, nil}

    for i := range testData {
        node := it.Next()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestDFSIteratorSequentialWithCycleGraph(t *testing.T) {
    g := NewGraph("g")
    a := g.NewSubGraph("a")
    b := g.NewSubGraph("b")
    c := g.NewSubGraph("c")

    a.ConnectNeighbour(b)
    a.ConnectNeighbour(c)
    b.ConnectNeighbour(c)
    b.ConnectNeighbour(a)
    c.ConnectNeighbour(a)
    c.ConnectNeighbour(b)

    it := a.NewDFSIterator()

    testData := []*Node{c, b, a, nil}

    for i := range testData {
        node := it.Next()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestDFSIteratorSimpleCircular(t *testing.T) {
    g := NewGraph("g")
    a := g.NewSubGraph("a")
    b := a.NewMutualNeighbour("b")
    it := a.NewDFSIterator()

    testData := []*Node{b, a, nil}

    for i := range testData {
        node := it.Next()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
    it = b.NewDFSIterator()
    testData = []*Node{a, b, nil}

    for i := range testData {
        node := it.Next()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func TestDFSIteratorSingleNode(t *testing.T) {
    g := NewGraph("g")
    a := g.NewSubGraph("a")
    it := a.NewDFSIterator()
    testData := []*Node{a, nil}

    for i := range testData {
        node := it.Next()
        if testData[i] != node {
            t.Error("DFS expected to deliver", testData[i], "but instead", node)
        }
    }
}

func createFullyConnectedGraph(numberOfNodes int) (*Node, *Node) {
    hypergraph := NewGraph("g")
    var subgraph *Node

    nodes := NodeSet(nil)

    for nodeid := 0; nodeid < numberOfNodes; nodeid++ {
        subgraph = hypergraph.NewSubGraph(fmt.Sprintf("%8d", nodeid))
        for _, previous := range nodes {
            subgraph.ConnectMutualNeighbour(previous)
        }
        nodes = append(nodes, subgraph)
    }
    return hypergraph, subgraph
}

func BenchmarkFullyConnectedDFSIterator(b *testing.B) {
    size := 1900
    fmt.Println(size)
    _, graph := createFullyConnectedGraph(size)
    b.ResetTimer()

    //fmt.Println("starting with node", graph, "\n")

    for i := 0; i < b.N; i++ {
        it := graph.NewDFSIterator()
        for node := it.Next(); node != nil; node = it.Next() {
           // fmt.Println(node)
        }
    }
}

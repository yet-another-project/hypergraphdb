package element_test
import (
    "testing"
    "github.com/yet-another-project/hypergraphdb/element"
)

func TestNewGraph(t *testing.T) {
    g := element.NewGraph("test")
    if "test" != g.String() {
        t.Error("wrong value")
    }
}

func TestNewSubGraph(t *testing.T) {
    g := element.NewGraph("g")
    h := g.NewSubGraph("h")
    if "g[h]" != g.String() {
        t.Error("wrong value")
    }
    if h != g.Subnodes()[0] {
        t.Error("wrong subnode")
    }
}

func TestNewNeighbour(t *testing.T) {
    g := element.NewGraph("g")
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
    if n != m.Neighbours()[0] {
        t.Error("wrong value")
    }
    if len(n.Neighbours()) != 0 {
        t.Error("n should not have neighbours")
    }
}

func TestNewMutualNeighbour(t *testing.T) {
    g := element.NewGraph("g")
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
    if n != m.Neighbours()[0] {
        t.Error("wrong value")
    }
    if len(n.Neighbours()) != 1 {
        t.Error("n should have neighbours")
    }
    if n.Neighbours()[0] != m {
        t.Error("m should also be n's neighbour")
    }

    if nil != g.NewMutualNeighbour("z") {
        t.Error("expected nil")
    }
}

func TestConnectNeighbour(t *testing.T) {
    g := element.NewGraph("g")
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
    a := element.NewGraph("a")
    b := a.NewSubGraph("b")
    c := b.NewSubGraph("c")
    d := c.NewSubGraph("d")
    expected := element.NodeSet(nil)
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
    a := element.NewGraph("a")
    b := a.NewSubGraph("b")
    c := b.NewSubGraph("c")
    d := c.NewSubGraph("d")
    e := b.NewSubGraph("e")

    testData := [][]*element.Node{
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


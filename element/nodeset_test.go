package element_test
import (
    "testing"
    "github.com/yet-another-project/hypergraphdb/element"
)

func TestFirstNodeNotIn(t *testing.T) {
    t.Log("TestFirstNodeNotIn")
    x := element.NewGraph("x")
    y := element.NewGraph("y")
    z := element.NewGraph("z")
    a := element.NewGraph("a")
    b := element.NewGraph("b")
    c := element.NewGraph("c")
    d := element.NewGraph("d")
    e := element.NewGraph("e")
    f := element.NewGraph("f")
    nodes_search := element.NewNodeSet(x, y, z)
    set1 := element.NewNodeSet(a, b, c)
    set2 := element.NewNodeSet(d, e, f)

    if x != nodes_search.FirstNodeNotIn(set2, set1) {
        t.Error("expected x")
    }
}

func TestFirstNodeNotInSecond(t *testing.T) {
    t.Log("TestFirstNodeNotInSecond")
    x := element.NewGraph("x")
    y := element.NewGraph("y")
    z := element.NewGraph("z")
    a := element.NewGraph("a")
    b := element.NewGraph("b")
    c := element.NewGraph("c")
    d := element.NewGraph("d")
    e := element.NewGraph("e")
    nodes_search := element.NewNodeSet(x, y, z)
    set1 := element.NewNodeSet(a, b, c)
    set2 := element.NewNodeSet(d, e, x)

    got := nodes_search.FirstNodeNotIn(set2, set1)
    if y != got {
        t.Error("expected y, got", got)
    }
}

func TestFirstNodeNotInAny(t *testing.T) {
    t.Log("TestFirstNodeNotInAny")
    a := element.NewGraph("a")
    b := element.NewGraph("b")
    c := element.NewGraph("c")
    d := element.NewGraph("d")
    e := element.NewGraph("e")
    f := element.NewGraph("f")
    nodes_search := element.NewNodeSet(a, e, c)
    set1 := element.NewNodeSet(a, b, c)
    set2 := element.NewNodeSet(d, e, f)

    got := nodes_search.FirstNodeNotIn(set2, set1)
    if nil != got {
        t.Error("expected nil, got", got)
    }
}

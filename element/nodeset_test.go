package element_test
import (
    "testing"
    "github.com/yet-another-project/hypergraphdb/element"
)

func TestFirstNodeNotIn(t *testing.T) {
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

func TestNodeSetCommonAncestorSimple(t *testing.T) {
    g := element.NewGraph("g")
    a := g.NewSubGraph("a")
    b := a.NewMutualNeighbour("b")
    c := a.NewMutualNeighbour("c")
    d := a.NewMutualNeighbour("d")
    set := element.NewNodeSet(b, c, d)
    ancestor := set.CommonAncestor()
    if g != ancestor {
        t.Error("expected common ancestor", g, "got", ancestor)
    }
}

func TestNodeSetCommonAncestorTree(t *testing.T) {
    g := element.NewGraph("g")
    g.ShowNeighbours = false
    g.ShowSubnodes = false
    one := g.NewSubGraph("1")
    two := one.NewSubGraph("2")
    two.NewSubGraph("3")
    zero := two.NewMutualNeighbour("0")
    four := zero.NewSubGraph("4")
    five := four.NewMutualNeighbour("5")
    b := four.NewSubGraph("b")
    six := b.NewMutualNeighbour("6")
    seven := five.NewSubGraph("7")
    eight := six.NewSubGraph("8")
    d := eight.NewMutualNeighbour("d")
    c := seven.NewSubGraph("c")
    a := eight.NewSubGraph("a")

    set := element.NewNodeSet(a, b, c, d)
    ancestor := set.CommonAncestor()
    if zero != ancestor {
        t.Error("expected common ancestor", zero, "got", ancestor)
    }
}

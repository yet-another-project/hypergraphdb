package graphdb
import (
    "testing"
)

func TestFirstNodeNotIn(t *testing.T) {
    x := NewGraph("x")
    y := NewGraph("y")
    z := NewGraph("z")
    a := NewGraph("a")
    b := NewGraph("b")
    c := NewGraph("c")
    d := NewGraph("e")
    e := NewGraph("f")
    f := NewGraph("f")
    nodes_search := NewNodeSet(x, y, z)
    set1 := NewNodeSet(a, b, c)
    set2 := NewNodeSet(d, e, f)

    if x != nodes_search.FirstNodeNotIn(set2, set1) {
        t.Error("expected x")
    }
}

package iterator

import (
    "github.com/yet-another-project/hypergraphdb/element"
)

type I interface {
    Stream() <-chan *element.Node
    Close()
    Next() *element.Node
    Run()
}

type iteratorContext struct {
    contextNode *element.Node
    neighbourIndex int
}

type iteratorContextStack []*iteratorContext

func (ctx *iteratorContextStack) TopNode() *element.Node {
    if len(*ctx) == 0 {
        return nil
    }
    return (*ctx)[len(*ctx)-1].contextNode
}

func (ctx *iteratorContextStack) NodeSet() element.NodeSet {
    lst := element.NodeSet(nil)
    for _, ctx2 := range *ctx {
        lst = append(lst, ctx2.contextNode)
    }
    return lst
}

func (ctx *iteratorContextStack) String() string {
    lst := ctx.NodeSet()
    return lst.String()
}

func (ctx *iteratorContextStack) AdvanceNeighbour() {
    if len(*ctx) > 0 {
        (*ctx)[len(*ctx)-1].neighbourIndex++
    }
}

func (ctx *iteratorContextStack) PopNode() *element.Node {
    if len(*ctx) == 0 {
        return nil
    }
    var node *element.Node
    node, *ctx = (*ctx)[len(*ctx)-1].contextNode, (*ctx)[:len(*ctx)-1]
    return node
}

func (ctx *iteratorContextStack) CurrentNeighbourIndex() int {
    if len(*ctx) > 0 {
        return (*ctx)[len(*ctx)-1].neighbourIndex
    }
    return -1
}

func (ctx *iteratorContextStack) HasMoreNeighbours() bool {
    if len(*ctx) > 0 && len(ctx.TopNode().Neighbours()) > (*ctx)[len(*ctx)-1].neighbourIndex {
        return true
    }
    return false
}

func (ctx *iteratorContextStack) CurrentNeighbour() *element.Node {
    if len(*ctx) == 0 {
        return nil
    }
    idx := (*ctx)[len(*ctx)-1].neighbourIndex
    if idx < 0 || idx >= len((*ctx)[len(*ctx)-1].contextNode.Neighbours()) {
        return nil
    }
    return (*ctx)[len(*ctx)-1].contextNode.Neighbours()[idx]
}

func (ctx *iteratorContextStack) NextNeighbour() *element.Node {
    if len(*ctx) == 0 {
        return nil
    }
    ctxFrame := (*ctx)[len(*ctx)-1]
    idx := ctxFrame.neighbourIndex + 1
    if idx < 0 || idx >= len(ctxFrame.contextNode.Neighbours()) {
        return nil
    }
    return ctxFrame.contextNode.Neighbours()[idx]
}

func (ctx *iteratorContextStack) PushNeighbour() {
    newctx := &iteratorContext{
        ctx.CurrentNeighbour(),
        -1,
    }
    *ctx = append(*ctx, newctx)
}

func (ctx *iteratorContextStack) PushNode(node *element.Node) *iteratorContext {
    newctx := &iteratorContext{
        node,
        -1,
    }
    *ctx = append(*ctx, newctx)
    return newctx
}


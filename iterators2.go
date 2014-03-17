package graphdb

import (
    "time"
//    "fmt"
//    "github.com/golang/glog"
)

type Iterator interface {
    Stream() <-chan *Node
    Close()
}

type DFSIterator struct {
    stream chan *Node
    closing chan bool
    visited NodeSet
    contextStack iteratorContextStack
}

type iteratorContext struct {
    contextNode *Node
    neighbourIndex int
}

type iteratorContextStack []*iteratorContext

func (n *Node) NewDFSIterator() *DFSIterator {
    it := &DFSIterator{
        stream: make(chan *Node),
        closing: make(chan bool),
        visited: NodeSet(nil),
        contextStack: iteratorContextStack(nil),
    }
    it.contextStack = append(it.contextStack, &iteratorContext{contextNode: n, neighbourIndex: -1})
    return it
}

func (it *DFSIterator) Stream() <-chan *Node {
    return nil
}

func (it *DFSIterator) Close() {
    it.closing<-true
    close(it.closing)
}

func (it *DFSIterator) returnNode(node *Node) *Node {
    alreadyAdded := false
    for _, visited := range it.visited {
        if visited == node {
            alreadyAdded = true
        }
    }
    if !alreadyAdded {//TODO: use instead a recursivity depth check and only append it when 0
        it.visited = append(it.visited, node)
    }
    return node
}

func (it *DFSIterator) hasBeenVisited(node *Node) bool {
    for _, visited := range it.visited {
        if visited == node {
            return true
        }
    }
    return false
}


func (ctx *iteratorContextStack) TopNode() *Node {
    if len(*ctx) == 0 {
        return nil
    }
    return (*ctx)[len(*ctx)-1].contextNode
}

func (ctx *iteratorContextStack) NodeSet() NodeSet {
    lst := NodeSet(nil)
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

func (ctx *iteratorContextStack) PopNode() *Node {
    if len(*ctx) == 0 {
        return nil
    }
    var node *Node
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
    if len(*ctx) > 0 && len(ctx.TopNode().neighbours) > (*ctx)[len(*ctx)-1].neighbourIndex {
        return true
    }
    return false
}

func (ctx *iteratorContextStack) CurrentNeighbour() *Node {
    if len(*ctx) == 0 {
        return nil
    }
    idx := (*ctx)[len(*ctx)-1].neighbourIndex
    if idx < 0 || idx >= len((*ctx)[len(*ctx)-1].contextNode.neighbours) {
        return nil
    }
    return (*ctx)[len(*ctx)-1].contextNode.neighbours[idx]
}

func (ctx *iteratorContextStack) NextNeighbour() *Node {
    if len(*ctx) == 0 {
        return nil
    }
    ctxFrame := (*ctx)[len(*ctx)-1]
    idx := ctxFrame.neighbourIndex + 1
    if idx < 0 || idx >= len(ctxFrame.contextNode.neighbours) {
        return nil
    }
    return ctxFrame.contextNode.neighbours[idx]
}

func (ctx *iteratorContextStack) PushNeighbour() {
    newctx := &iteratorContext{
        ctx.CurrentNeighbour(),
        -1,
    }
    *ctx = append(*ctx, newctx)
}

func (ctx *iteratorContextStack) PushNode(node *Node) *iteratorContext {
    newctx := &iteratorContext{
        node,
        -1,
    }
    *ctx = append(*ctx, newctx)
    return newctx
}

func (it *DFSIterator) Next() *Node {
    //vlevel := glog.Level(1)
    //glog.V(vlevel).Infoln("\n------")
    topNode := it.contextStack.TopNode()
    //glog.V(vlevel).Infoln("top node", topNode)
    //glog.V(vlevel).Infoln("visited", it.visited.String())
    //glog.V(vlevel).Infoln("stack", it.contextStack.String(), "\n-")
    if nil == topNode {
        return nil
    }

    nextNeighbour := topNode.neighbours.FirstNodeNotIn(it.visited, it.contextStack.NodeSet())
    //glog.V(vlevel).Infoln("NEIGHBOURS next", nextNeighbour)

    if len(it.contextStack) == 0 {
        //glog.V(vlevel).Infoln("RET nil")
        return nil
    }
    if nextNeighbour == nil {
        //glog.V(vlevel).Infoln("nextNeighbour nil, returning", it.contextStack.TopNode())
        popped := it.contextStack.PopNode()
        //glog.V(vlevel).Infoln("popped", popped, "top", topNode)
        return it.returnNode(popped)
    }
    if !it.visited.ContainsNode(nextNeighbour) {
        it.contextStack.PushNode(nextNeighbour)
        //glog.V(vlevel).Infoln("pushed, new stack", it.contextStack.String(), "\n-")
        return it.returnNode(it.Next())
    } else {
        panic("should never happen")
    }

    return nil
}

func (it *DFSIterator) Run() {
    var pending NodeSet
    for {
        startFetch := time.After(time.Duration(0))
        var first *Node
        var stream chan *Node
        if len(pending) > 0 {
            first = pending[0]
            stream = it.stream
        }
        select {
        case <- startFetch:
            next := it.Next()
            if nil == next {
                close(it.stream)
                it.Close()
            }
            pending = append(pending, next)
        case <- it.closing:
            close(it.stream)
            return
        case stream<- first:
            pending = pending[1:]
        }
    }
}

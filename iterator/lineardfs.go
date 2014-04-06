package iterator

import (
    "time"
    "github.com/yet-another-project/hypergraphdb/element"
)

type LinearDFS struct {
    stream chan *element.Node
    closing chan bool
    visited element.NodeSet
    contextStack iteratorContextStack
}

func NewLinearDFS(n *element.Node) *LinearDFS {
    it := &LinearDFS{
        stream: make(chan *element.Node),
        closing: make(chan bool),
        visited: element.NodeSet(nil),
        contextStack: iteratorContextStack(nil),
    }
    it.contextStack = append(it.contextStack, &iteratorContext{contextNode: n, neighbourIndex: -1})
    return it
}

func (it *LinearDFS) Stream() <-chan *element.Node {
    return it.stream
}

func (it *LinearDFS) Close() {
    it.closing<-true
    close(it.closing)
}

func (it *LinearDFS) returnNode(node *element.Node) *element.Node {
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

func (it *LinearDFS) hasBeenVisited(node *element.Node) bool {
    for _, visited := range it.visited {
        if visited == node {
            return true
        }
    }
    return false
}

func (it *LinearDFS) Next() *element.Node {
    //vlevel := glog.Level(1)
    //glog.V(vlevel).Infoln("\n------")
    topNode := it.contextStack.TopNode()
    //glog.V(vlevel).Infoln("top node", topNode)
    //glog.V(vlevel).Infoln("visited", it.visited.String())
    //glog.V(vlevel).Infoln("stack", it.contextStack.String(), "\n-")
    if nil == topNode {
        return nil
    }

    nextNeighbour := topNode.Neighbours().FirstNodeNotIn(it.visited, it.contextStack.NodeSet())
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
    if _, ok := it.visited.ContainsNode(nextNeighbour); !ok {
        it.contextStack.PushNode(nextNeighbour)
        //glog.V(vlevel).Infoln("pushed, new stack", it.contextStack.String(), "\n-")
        return it.returnNode(it.Next())
    } else {
        panic("should never happen")
    }

    return nil
}

func (it *LinearDFS) Run() {
    var pending element.NodeSet
    for {
        startFetch := time.After(time.Duration(0))
        var first *element.Node
        var stream chan *element.Node
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


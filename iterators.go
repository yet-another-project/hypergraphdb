package graphdb

import (
    "github.com/golang/glog"
)

type Iterator interface {
    Stream() <-chan *Node
}

type ChildrenIterator struct {
    start *Node
}

func NewChildrenIterator(start *Node) *ChildrenIterator {
    return &ChildrenIterator{start}
}

func (it *ChildrenIterator) Stream() <-chan *Node {
    ch := make(chan *Node)
    go func() {
        for child := range it.start.subnodes {
            ch <- it.start.subnodes[child]
        }
        close(ch)
    }()
    return ch
}

type UpwardParentIterator struct {
    start *Node
}

func (it *UpwardParentIterator) Stream() <-chan *Node {
    ch := make(chan *Node)
    currentNode := it.start.parent
    go func() {
        for currentNode != nil {
            ch <- currentNode
            currentNode = currentNode.parent
        }
        close(ch)
    }()
    return ch
}

func NewUpwardParentIterator(start *Node) *UpwardParentIterator {
    return &UpwardParentIterator{start}
}

type DFSIterator struct {
    start *Node
    current NodeSet
    visited NodeSet
    stream chan *Node
}

func NewDFSIterator(start *Node) *DFSIterator {
    it := &DFSIterator{start, NodeSet(nil), NodeSet(nil), make(chan *Node)}
    it.current = append(it.current, start)
    return it
}

func (it *DFSIterator) dfsUtil(pre bool) {
    stackTop := it.current[len(it.current)-1]
    glog.V(1).Infoln("started ctx " + stackTop.String())
    if !it.isVisited(stackTop) && pre {
        glog.V(2).Infoln("pushing at start of scope " + stackTop.String())
        it.pushNode(stackTop)
        status := it.popNode(stackTop)
        if !status {
            glog.V(2).Infoln("could not pop node")
        }
    }
    for _, neighbour := range stackTop.neighbours {
            if !it.isVisited(neighbour) && !it.isOnBackPath(neighbour) {
                it.current = append(it.current, neighbour)
                it.visited = append(it.visited, neighbour)
                if pre {
                    glog.V(2).Infoln("before recursive, pushing " + neighbour.String() + " from ctx " + stackTop.String())
                    it.pushNode(neighbour)
                }
                it.dfsUtil(pre)
                if !pre {
                    glog.V(2).Infoln("back from recursive, pushing " + neighbour.String() + " from ctx " + stackTop.String())
                    it.pushNode(neighbour)
                }
            }
    }
    if !it.isVisited(stackTop) && !pre {
        glog.V(2).Infoln("pushing before end of scope " + stackTop.String())
        it.pushNode(stackTop)
        status := it.popNode(stackTop)
        if !status {
            glog.V(2).Infoln("could not pop node")
        }
    }
}

func (it *DFSIterator) Stream() <-chan *Node {
    go func() {
        it.dfsUtil(false)
        close(it.stream)
    }()
    return it.stream
}

func (it *DFSIterator) isVisited(node *Node) bool {
    for _, visited := range it.visited {
        if visited == node {
            return true
        }
    }
    return false
}

func (it *DFSIterator) isOnBackPath(node *Node) bool {
    for _, prev := range it.current {
        if prev == node {
            return true
        }
    }
    return false
}

func (it *DFSIterator) pushNode(node *Node) bool {
    glog.V(1).Infoln("pushing node", node)
    it.stream<- node
    it.visited = append(it.visited, node)
    return true
}

func (it *DFSIterator) popNode(node *Node) bool {
    head, tail := it.visited[len(it.visited)-1], it.visited[:len(it.visited)-1]
    if head == node {
        it.visited = tail
        return true
    }
    return false
}

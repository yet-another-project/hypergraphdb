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

func (self *ChildrenIterator) Stream() <-chan *Node {
    ch := make(chan *Node)
    go func() {
        for child := range self.start.subnodes {
            ch <- self.start.subnodes[child]
        }
        close(ch)
    }()
    return ch
}

type UpwardParentIterator struct {
    start *Node
}

func (self *UpwardParentIterator) Stream() <-chan *Node {
    ch := make(chan *Node)
    currentNode := self.start.parent
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

func (self *DFSIterator) dfsUtil(pre bool) {
    stackTop := self.current[len(self.current)-1]
    glog.V(1).Infoln("started ctx " + stackTop.String())
    if !self.isVisited(stackTop) && pre {
        glog.V(2).Infoln("pushing at start of scope " + stackTop.String())
        self.pushNode(stackTop)
        status := self.popNode(stackTop)
        if !status {
            glog.V(2).Infoln("could not pop node")
        }
    }
    for _, neighbour := range stackTop.neighbours {
            if !self.isVisited(neighbour) && !self.isOnBackPath(neighbour) {
                self.current = append(self.current, neighbour)
                self.visited = append(self.visited, neighbour)
                if pre {
                    glog.V(2).Infoln("before recursive, pushing " + neighbour.String() + " from ctx " + stackTop.String())
                    self.pushNode(neighbour)
                }
                self.dfsUtil(pre)
                if !pre {
                    glog.V(2).Infoln("back from recursive, pushing " + neighbour.String() + " from ctx " + stackTop.String())
                    self.pushNode(neighbour)
                }
            }
    }
    if !self.isVisited(stackTop) && !pre {
        glog.V(2).Infoln("pushing before end of scope " + stackTop.String())
        self.pushNode(stackTop)
        status := self.popNode(stackTop)
        if !status {
            glog.V(2).Infoln("could not pop node")
        }
    }
}

func (self *DFSIterator) Stream() <-chan *Node {
    go func() {
        self.dfsUtil(false)
        close(self.stream)
    }()
    return self.stream
}

func (self *DFSIterator) isVisited(node *Node) bool {
    for _, visited := range self.visited {
        if visited == node {
            return true
        }
    }
    return false
}

func (self *DFSIterator) isOnBackPath(node *Node) bool {
    for _, prev := range self.current {
        if prev == node {
            return true
        }
    }
    return false
}

func (self *DFSIterator) pushNode(node *Node) bool {
    glog.V(1).Infoln("pushing node", node)
    self.stream<- node
    self.visited = append(self.visited, node)
    return true
}

func (self *DFSIterator) popNode(node *Node) bool {
    head, tail := self.visited[len(self.visited)-1], self.visited[:len(self.visited)-1]
    if head == node {
        self.visited = tail
        return true
    }
    return false
}

package iterator
import (
    "github.com/yet-another-project/hypergraphdb/element"
    "github.com/golang/glog"
)

type RecursiveDFS struct {
    start *element.Node
    current element.NodeSet
    visited element.NodeSet
    stream chan *element.Node
}

func NewRecursiveDFS(start *element.Node) *RecursiveDFS {
    it := &RecursiveDFS{start, element.NodeSet(nil), element.NodeSet(nil), make(chan *element.Node)}
    it.current = append(it.current, start)
    return it
}

func (it *RecursiveDFS) dfsUtil(pre bool) {
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
    for _, neighbour := range stackTop.Neighbours() {
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

func (it *RecursiveDFS) Stream() <-chan *element.Node {
    go func() {
        it.dfsUtil(false)
        close(it.stream)
    }()
    return it.stream
}

func (it *RecursiveDFS) isVisited(node *element.Node) bool {
    for _, visited := range it.visited {
        if visited == node {
            return true
        }
    }
    return false
}

func (it *RecursiveDFS) isOnBackPath(node *element.Node) bool {
    for _, prev := range it.current {
        if prev == node {
            return true
        }
    }
    return false
}

func (it *RecursiveDFS) pushNode(node *element.Node) bool {
    glog.V(1).Infoln("pushing node", node)
    it.stream<- node
    it.visited = append(it.visited, node)
    return true
}

func (it *RecursiveDFS) popNode(node *element.Node) bool {
    head, tail := it.visited[len(it.visited)-1], it.visited[:len(it.visited)-1]
    if head == node {
        it.visited = tail
        return true
    }
    return false
}


package graphdb
import (
    //"github.com/golang/glog"
)

func (nodes *NodeSet) String() string {
    str := "["
    for i := range *nodes {
        str += (*nodes)[i].String()
        if i != len(*nodes)-1 {
            str += ", "
        }
    }
    return str + "]"
}

func (nodes NodeSet) Intersect(nodes2 NodeSet) NodeSet {
    counted := make(map[*Node]bool)
    for _, node := range nodes {
        counted[node] = true
    }
    common := NodeSet(nil)
    for _, node := range nodes2 {
        if _, ok := counted[node]; ok {
            common = append(common, node)
        }
    }
    return common
}

func (nodes NodeSet) Union(nodes2 NodeSet) NodeSet {
    return nil
}

func (nodes NodeSet) Difference(nodes2 NodeSet) NodeSet {
    return nil
}

func (nodes NodeSet) Xor(nodes2 NodeSet) NodeSet {
    return nil
}

func (nodes NodeSet) ContainsNode(node *Node) bool {
    for _, localNode := range nodes {
        if localNode == node {
            return true
        }
    }
    return false
}

func (nodes NodeSet) ContainsSubset(otherSet NodeSet) bool {
    return false
}

// TODO use a concurrent version
func (nodes NodeSet) FirstNodeNotIn(sets ...NodeSet) *Node {
    outerLoop:
    for _, localNode := range nodes {
        for _, set := range sets {
            if set.ContainsNode(localNode) {
                continue outerLoop
            }
        }
        return localNode
    }
    return nil
}

func NewNodeSet(nodes ...*Node) NodeSet {
    nodeset := NodeSet(nil)
    for _, node := range nodes {
        nodeset = append(nodeset, node)
    }
    return nodeset
}

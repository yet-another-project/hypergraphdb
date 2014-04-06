package element

import (
)

type NodeSet []*Node

func (set NodeSet) String() string {
    str := "["
    for i := range set {
        str += (set)[i].String()
        if i != len(set)-1 {
            str += ", "
        }
    }
    return str + "]"
}

func (set NodeSet) Intersect(nodes2 NodeSet) NodeSet {
    counted := make(map[*Node]bool)
    for _, node := range set {
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

func (set NodeSet) Union(nodes2 NodeSet) NodeSet {
    return nil
}

func (set NodeSet) Difference(nodes2 NodeSet) NodeSet {
    return nil
}

func (set NodeSet) Xor(nodes2 NodeSet) NodeSet {
    return nil
}

func (set NodeSet) ContainsNode(node *Node) (int, bool) {
    for position, localNode := range set {
        if localNode == node {
            return position, true
        }
    }
    return -1, false
}

func (set NodeSet) ContainsSubset(otherSet NodeSet) bool {
    return false
}

// TODO use a concurrent version
func (set NodeSet) FirstNodeNotIn(sets ...NodeSet) *Node {
    outerLoop:
    for _, localNode := range set {
        for _, set := range sets {
            if _, ok := set.ContainsNode(localNode); ok {
                continue outerLoop
            }
        }
        return localNode
    }
    return nil
}

func (set NodeSet) CommonAncestor() *Node {
    if len(set) == 0 {
        return nil
    }
    if len(set) == 1 {
        return set[0]
    }
    currentNode := set[0]
    pruningSubject := currentNode.UpwardParents()
    nodesLeft := set[1:]
    for _, currentNode := range nodesLeft {
        if pruneAtPos, ok := pruningSubject.ContainsNode(currentNode); ok {
            pruningSubject = append(pruningSubject[:pruneAtPos], pruningSubject[pruneAtPos+1:]...)
        } else {
            currentAncestors := currentNode.UpwardParents()
            pruningSubject = pruningSubject.Intersect(currentAncestors)
        }
    }
    if len(pruningSubject) > 0 {
        return pruningSubject[0]
    }
    return nil
}

func NewNodeSet(set ...*Node) NodeSet {
    nodeset := NodeSet(nil)
    for _, node := range set {
        nodeset = append(nodeset, node)
    }
    return nodeset
}

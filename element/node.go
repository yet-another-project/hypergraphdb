package element

import "github.com/golang/glog"

type Node struct {
    label string
    parent *Node
    subnodes NodeSet //nested subgraphs
    neighbours NodeSet //connected nodes (regular nodes, regular edges)
    hypertrail NodeSet //the nodes this node goes through as a hyperedge
}

func NewGraph(label string) *Node {
    return &Node{label, nil, NodeSet(nil), NodeSet(nil), NodeSet(nil)}
}

func (parent *Node) NewSubGraph(label string) *Node {
    newNode := NewGraph(label)
    newNode.parent = parent
    parent.subnodes = append(parent.subnodes, newNode)
    return newNode
}

func (node *Node) NewNeighbour(label string) *Node {
    if nil == node.parent {
        return nil
    }
    newNode := node.parent.NewSubGraph(label)
    node.neighbours = append(node.neighbours, newNode)
    return newNode
}

func (node *Node) NewMutualNeighbour(label string) *Node {
    newNode := node.NewNeighbour(label)
    if nil == newNode {
        return nil
    }
    newNode.ConnectNeighbour(node)
    return newNode
}

func (node *Node) ConnectNeighbour(other *Node) bool {
    for _, neighbour := range node.neighbours {
        if neighbour == other {
            glog.V(1).Infoln(other.String() + " is already a neighbour of " + node.String())
            return false
        }
    }
    node.neighbours = append(node.neighbours, other)
    return true
}

func (node *Node) ConnectMutualNeighbour(other *Node) bool {
    if !node.ConnectNeighbour(other) {
        return false
    }
    if !other.ConnectNeighbour(node) {
        //TODO disconnect
        return false
    }
    return true
}

//------------------- exploration
func (node *Node) UpwardParents() NodeSet {
    parents := NodeSet(nil)
    parent := node.parent
    for parent != nil {
        parents = append(parents, parent)
        parent = parent.parent
    }
    return parents
}

func (node *Node) CommonAncestor(other *Node) *Node {
    nodes1 := node.UpwardParents()
    nodes2 := other.UpwardParents()

    if len(nodes1) > len(nodes2) {
        nodes2, nodes1 = nodes1, nodes2
    }

    var common *Node
outerLoop:
    for _, node2 := range nodes2 {
        for _, node1 := range nodes1 {
            if node1 == node2 {
                common = node1
                break outerLoop
            }
        }
    }
    return common
}

//------------------- information
func (node *Node) Subnodes() NodeSet {
    return node.subnodes
}

func (node *Node) Neighbours() NodeSet {
    return node.neighbours
}

func (parent *Node) String() string {
    str := parent.label
    if len(parent.neighbours) > 0 {
        str += "{"
        for i := range parent.neighbours {
            str += parent.neighbours[i].label
            if i != len(parent.neighbours)-1 {
                str += ", "
            }
        }
        str += "}"
    }
    if len(parent.hypertrail) > 0 {
        if len(parent.neighbours) > 0 {
            str += "  "
        }
        str += "<"
        for i := range parent.hypertrail {
            str += parent.hypertrail[i].label
            if i != len(parent.hypertrail)-1 {
                str += ", "
            }
        }
        str += ">"
    }
    if len(parent.subnodes) > 0 {
        if len(parent.hypertrail) > 0 {
            str += "  "
        }
        str += "["
        for i := range parent.subnodes {
            str += parent.subnodes[i].label
            if i != len(parent.subnodes)-1 {
                str += ", "
            }
        }
        str += "]"
    }
    return str
}


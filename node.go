package graphdb
import (
    "fmt"
)

type Node struct {
    label string
    parent *Node
    subnodes NodeSet //nested subgraphs (hyperedges)
    neighbours NodeSet //connected nodes (regular nodes, regular edges)
}

type NodeSet []*Node

func NewGraph(label string) *Node {
    return &Node{label, nil, NodeSet(nil), NodeSet(nil)}
}

func (parent *Node) NewNode(label string) *Node {
    newNode := NewGraph(label)
    newNode.parent = parent
    parent.subnodes = append(parent.subnodes, newNode)
    return newNode
}

func (parent *Node) AttachNode(node *Node) {
    if node.parent != nil {
        for child := range node.parent.subnodes {
            fmt.Println(child)
        }
    }

    parent.subnodes = append(parent.subnodes, node)
    node.parent = parent
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
    if len(parent.subnodes) > 0 {
        if len(parent.neighbours) > 0 {
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

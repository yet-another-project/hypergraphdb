package iterator_test
import (
    "fmt"
    "github.com/yet-another-project/hypergraphdb/element"
)

func createFullyConnectedGraph(numberOfNodes int) (*element.Node, *element.Node) {
    hypergraph := element.NewGraph("g")
    var subgraph *element.Node

    nodes := element.NodeSet(nil)

    for nodeid := 0; nodeid < numberOfNodes; nodeid++ {
        subgraph = hypergraph.NewSubGraph(fmt.Sprintf("%8d", nodeid))
        for _, previous := range nodes {
            subgraph.ConnectMutualNeighbour(previous)
        }
        nodes = append(nodes, subgraph)
    }
    return hypergraph, subgraph
}

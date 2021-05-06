package util

import (
	dfa "github.com/skius/dataflowanalysis"
	"github.com/skius/stringlang/cfg"
	"sort"
)

// NodePI is a concrete implementation of the dfa.NodePI interface, useful for stringlang CFGs
type NodePI struct {
	inner *cfg.Node
}

func (n *NodePI) CFGNode() *cfg.Node {
	return n.inner
}

func (n *NodePI) Label() int {
	return n.inner.Label
}

func (n *NodePI) Preds() []int {
	preds := make([]int, 0, len(n.inner.PredsNotTaken)+len(n.inner.PredsTaken))
	for _, p := range n.inner.PredsNotTaken {
		preds = append(preds, p.Label)
	}
	for _, p := range n.inner.PredsTaken {
		preds = append(preds, p.Label)
	}
	return preds
}

func (n *NodePI) Succs() []int {
	succs := make([]int, 0, 2)
	if n.inner.SuccNotTaken != nil {
		succs = append(succs, n.inner.SuccNotTaken.Label)
	}

	if n.inner.SuccTaken != nil {
		succs = append(succs, n.inner.SuccTaken.Label)
	}

	return succs
}

func (n *NodePI) Get() dfa.Stmt {
	return n.inner.Expr
}

// NodePIFromCFG returns the slice of all labels and the label-to-NodePI-interface map of the input CFG
func NodePIFromCFG(graph *cfg.CFG) ([]int, map[int]dfa.NodePI) {
	ids := make([]int, 0)
	idToNode := make(map[int]dfa.NodePI)

	graph.Visit(func(node *cfg.Node) {
		ids = append(ids, node.Label)

		dfaNode := new(NodePI)
		dfaNode.inner = node
		idToNode[node.Label] = dfaNode
	})

	sort.Ints(ids)
	return ids, idToNode
}

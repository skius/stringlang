package util

import (
	dfa "github.com/skius/dataflowanalysis"
	"github.com/skius/stringlang/cfg"
	"sort"
)

// NodePS is a concrete implementation of the path-sensitive dataflowanalysis.Node interface
type NodePS struct {
	inner *cfg.Node
}

func (n *NodePS) Label() int {
	return n.inner.Label
}

func (n *NodePS) PredsNotTaken() []int {
	preds := make([]int, len(n.inner.PredsNotTaken))
	for i, p := range n.inner.PredsNotTaken {
		preds[i] = p.Label
	}
	return preds
}

func (n *NodePS) PredsTaken() []int {
	preds := make([]int, len(n.inner.PredsTaken))
	for i, p := range n.inner.PredsTaken {
		preds[i] = p.Label
	}
	return preds
}

func (n *NodePS) SuccsNotTaken() []int {
	if n.inner.SuccNotTaken == nil {
		return []int{}
	}

	return []int{n.inner.SuccNotTaken.Label}
}

func (n *NodePS) SuccsTaken() []int {
	if n.inner.SuccTaken == nil {
		return []int{}
	}

	return []int{n.inner.SuccTaken.Label}
}

func (n *NodePS) Get() dfa.Stmt {
	return n.inner.Expr
}

// NodePSFromCFG returns the slice of all labels and the label-to-NodePS-interface map of the input CFG
func NodePSFromCFG(graph *cfg.CFG) ([]int, map[int]dfa.Node) {
	ids := make([]int, 0)
	idToNode := make(map[int]dfa.Node)

	graph.Visit(func(node *cfg.Node) {
		ids = append(ids, node.Label)

		dfaNode := new(NodePS)
		dfaNode.inner = node
		idToNode[node.Label] = dfaNode
	})

	sort.Ints(ids)
	return ids, idToNode
}

package cfg

import (
	"github.com/skius/stringlang/ast"
	"math"
)

// ExprWithNode is an Expr wrapper used to still be able to use the analyses' results (which are mappings from
// Node.Label to Facts) in Expr form, instead of CFG form
type ExprWithNode struct {
	ast.Expr
	Node *Node
}

func (cfg *CFG) Collect() ast.Block {
	block := collectBlock(cfg.Entry, []int{}, []ast.Expr{})
	return block
}

func collectBlock(n *Node, stopList []int, prev []ast.Expr) []ast.Expr {
	if n == nil {
		return prev
	}
	if n.IsSentinel() {
		return collectBlock(n.SuccNotTaken, stopList, prev)
	}
	for _, l := range stopList {
		// If this is a stop node, we're already one too far and don't need to add n.Expr to prev
		if n.Label == l {
			return prev
		}
	}
	if !n.IsBranch() {
		//return collectBlock(n.SuccNotTaken, stopList, append(prev, ExprWithNode{Expr: n.Expr, Node: n}))
		return collectBlock(n.SuccNotTaken, stopList, append(prev, n.Expr))
	}
	if n.IsWhile() {
		wBlock := collectBlock(n.SuccTaken, append(stopList, n.Label), []ast.Expr{})
		w := ast.While{
			Cond: n.Expr,
			Body: ast.Block(wBlock),
		}
		//return collectBlock(n.SuccNotTaken, stopList, append(prev, ExprWithNode{Expr: w, Node: n}))
		return collectBlock(n.SuccNotTaken, stopList, append(prev, w))
	}
	if n.IsIf() {
		// Need to find merge point, i.e. node that is a descendant of the NT and also T branch
		reach := make([]int, 0)
		VisitAll(n.SuccNotTaken, func(node *Node) {
			reach = append(reach, node.Label)
		})
		merge := math.MaxInt32
		var mergeNode *Node
		VisitAll(n.SuccTaken, func(node *Node) {
			for _, r := range reach {
				if r == node.Label && r < merge {
					merge = r
					mergeNode = node
				}
			}
		})

		// If we don't have a merge point it means both branches lead to exits, and in that case the way we handle
		// nil works out

		then := collectBlock(n.SuccTaken, append(stopList, merge), []ast.Expr{})
		eelse := collectBlock(n.SuccNotTaken, append(stopList, merge), []ast.Expr{})
		ifelse := ast.IfElse{
			Cond: n.Expr,
			Then: ast.Block(then),
			Else: ast.Block(eelse),
		}
		//return collectBlock(mergeNode, stopList, append(prev, ExprWithNode{Expr: ifelse, Node: n}))
		return collectBlock(mergeNode, stopList, append(prev, ifelse))
	}

	return prev
}

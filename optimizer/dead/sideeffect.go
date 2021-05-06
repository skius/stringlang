package dead

import (
	"fmt"
	"github.com/skius/stringlang/ast"
	"github.com/skius/stringlang/cfg"
	"github.com/skius/stringlang/optimizer/analysis/sideeffect"
	"github.com/skius/stringlang/optimizer/analysis/util"
)

// Modifies g
func sideeffectDCE(g *cfg.CFG) *cfg.CFG {
	in, out := sideeffect.Compute(g)
	fmt.Println(util.PrettyPrintFlows(g, in, out))
	fmt.Println()

	//gv := cfg.GraphViz(g, nil)
	//fmt.Println(gv.String())

	g.Visit(func(node *cfg.Node) {
		// If it's an exit, don't remove it
		for _, exit := range g.Exits {
			if exit.Label == node.Label {
				return
			}
		}

		expr := node.Expr
		if assn, ok := expr.(ast.Assn); ok {
			if ast.HasSideEffects(assn.E) {
				return
			}
			// If the assignment defines no side-effect-live variables, we can remove it
			seLive := out[node.Label]
			if !seLive.Contains(string(assn.V)) {
				node.Remove()
			}
			//fmt.Println(expr.String())

			return
		} else {
			//fmt.Println(expr.String())
		}

		// You can remove an IfElse branch if both its successors point to the same node
		if node.IsIf() {
			if ast.HasSideEffects(expr) {
				return
			}
			if node.SuccNotTaken.Label == node.SuccTaken.Label {
				node.Remove()
			}
			return
		}

		// TODO: Handle Whiles
		if node.IsWhile() {
			return
		}

		// If it's not an assignment, not a branch and doesn't have side-effects, remove it too
		if !ast.HasSideEffects(expr) {
			node.Remove()
			return
		}

	})

	return g
}

func sideeffectDCEWhile(g *cfg.CFG) *cfg.CFG {
	in, out := sideeffect.Compute(g)
	fmt.Println("Flows in DCEWHILE")
	fmt.Println(util.PrettyPrintFlows(g, in, out))
	fmt.Println()
	collected := g.Collect()
	fmt.Println("Collected: ")
	fmt.Println(collected)
	fmt.Println()

	// Handle while-loop dead code elimination
	newBlock := make(ast.Block, 0, len(collected))
	for i := range collected {
		exprNE := collected[i]
		exprN := exprNE.(cfg.ExprWithNode)
		expr := exprN.Node.Expr
		fmt.Println("Looking at", expr.String())
		if _, ok := exprN.Expr.(ast.While); ok {
			if exprN.Node.SuccNotTaken == nil {
				newBlock = append(newBlock, expr)
				continue
			}
			seNT := in[exprN.Node.SuccNotTaken.Label]

			// If the while loop has side-effects, or defines a side-effect-live variable for its NT branch (seNT)
			// we can't optimize it away
			if ast.HasSideEffects(expr) {
				fmt.Println("has se", expr)
				newBlock = append(newBlock, expr)
				continue
			}
			defs := ast.DefinedVars(expr)
			for def := range defs {
				if seNT.Contains(def) {
					newBlock = append(newBlock, expr)
					continue
				}
			}

			// Otherwise we can optimize it away
			continue
		}
		newBlock = append(newBlock, expr)
	}

	return cfg.NewFromBlock(newBlock)
}

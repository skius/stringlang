package liveness

import (
	dfa "github.com/skius/dataflowanalysis"
	"github.com/skius/stringlang/ast"
	"github.com/skius/stringlang/cfg"
	"github.com/skius/stringlang/optimizer/analysis/util"
)

// Compute takes a CFG of a *normalized* program and computes the live_in and live_out sets of variables at each node
func Compute(graph *cfg.CFG) (liveIn, liveOut map[int]util.Set) {
	// Fill structures
	ids, idToNode := util.NodePIFromCFG(graph)

	bottom := util.EmptySet()

	// Liveness is a may analysis, hence the data-flow Meet operator is Union
	merge := func(s1F, s2F dfa.Fact) dfa.Fact {
		s1 := s1F.(util.Set)
		s2 := s2F.(util.Set)
		return s1.Union(s2)
	}

	// Liveness is a backward analysis, hence this function takes the live_out of a node and computes the live_in
	flow := func(setF dfa.Fact, nodeF dfa.NodePI) (res dfa.Fact) {
		set := setF.(util.Set)
		node := nodeF.(*util.NodePI)
		expr := node.Get().(ast.Expr)

		gen := util.Set(ast.UsedVars(expr))
		kill := make(util.Set)

		if val, ok := expr.(ast.Assn); ok {
			kill = util.SetFrom(string(val.V))
		}

		// inFlow = gen(node) \union (outFlow \except kill(node))
		res = gen.Union(set.Except(kill))

		return res
	}

	in, out := dfa.RunBackwardPI(ids, idToNode, merge, flow, bottom)

	liveIn, liveOut = make(map[int]util.Set), make(map[int]util.Set)

	for k, v := range in {
		liveIn[k] = v.(util.Set)
	}

	for k, v := range out {
		liveOut[k] = v.(util.Set)
	}

	return liveIn, liveOut
}

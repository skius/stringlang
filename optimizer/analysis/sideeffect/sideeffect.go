package sideeffect

import (
	dfa "github.com/skius/dataflowanalysis"
	"github.com/skius/stringlang/ast"
	"github.com/skius/stringlang/cfg"
	"github.com/skius/stringlang/optimizer/analysis/util"
)

// Compute takes a CFG of a *normalized* program and computes the in and out sets of variables at each node that
// are used for side-effects further on. Side-effects include: Used in the return value and used as arguments in a call.
func Compute(graph *cfg.CFG) (seLiveIn, seLiveOut map[int]util.Set) {
	exitLabels := make(map[int]struct{})
	for _, v := range graph.Exits {
		exitLabels[v.Label] = struct{}{}
	}

	// Fill structures
	ids, idToNode := util.NodePIFromCFG(graph)

	bottom := util.EmptySet()

	// May analysis, hence the data-flow Meet operator is Union
	merge := func(s1F, s2F dfa.Fact) dfa.Fact {
		s1 := s1F.(util.Set)
		s2 := s2F.(util.Set)
		return s1.Union(s2)
	}

	// Backward analysis, hence this function takes the seLiveOut of a node and computes the seLiveIn
	flow := func(setF dfa.Fact, nodeF dfa.NodePI) (res dfa.Fact) {
		set := setF.(util.Set)
		node := nodeF.(*util.NodePI)
		expr := node.Get().(ast.Expr)

		if _, ok := exitLabels[node.Label()]; ok {
			// Current node is n exit node, i.e. a return value. All UsedVars are therefore used in for side-effects
			// Additionally we don't have any useful incoming facts, hence we can just return
			return util.Set(ast.UsedVars([]ast.Expr{expr}))
		}

		gen := make(util.Set)
		kill := make(util.Set)

		if val, ok := expr.(ast.Assn); ok {
			// Same as for liveness, if we define a variable it will not be side-effect-live before it's defined
			kill = util.SetFrom(string(val.V))

			if set.Contains(string(val.V)) {
				// Additionally however, the variables used to define side-effect-live variables are now
				// side-effect-live
				gen = ast.UsedVars([]ast.Expr{val.E})
			}
		}

		if val, ok := expr.(ast.Call); ok {
			// A call may contain side-effects, hence all variables used as arguments are side-effect-live
			for _, arg := range val.Args {
				gen = gen.Union(ast.UsedVars([]ast.Expr{arg}))
			}
		}

		// inFlow = gen(node) \union (outFlow \except kill(node))
		res = gen.Union(set.Except(kill))

		return res
	}

	in, out := dfa.RunBackwardPI(ids, idToNode, merge, flow, bottom)

	seLiveIn, seLiveOut = make(map[int]util.Set), make(map[int]util.Set)

	for k, v := range in {
		seLiveIn[k] = v.(util.Set)
	}

	for k, v := range out {
		seLiveOut[k] = v.(util.Set)
	}

	return seLiveIn, seLiveOut
}

package util

import (
	dfa "github.com/skius/dataflowanalysis"
	"github.com/skius/stringlang/cfg"
	"sort"
	"strings"
)

// A Set is a dataflowanalysis.Fact compliant lattice of sets (e.g. the lattice of the powerset of all variables)
type Set map[string]struct{}

func (s Set) Union(s2 Set) Set {
	u := make(Set, len(s))
	for k, v := range s {
		u[k] = v
	}

	for k, v := range s2 {
		u[k] = v
	}

	return u
}

func (s Set) Intersect(s2 Set) Set {
	u := make(Set, len(s))
	for k, v := range s {
		if _, ok := s2[k]; ok {
			u[k] = v
		}
	}

	return u
}

func (s Set) Except(s2 Set) Set {
	u := make(Set, len(s))
	for k, v := range s {
		u[k] = v
	}

	for k := range s2 {
		delete(u, k)
	}

	return u
}

func (s Set) Contains(v string) bool {
	_, ok := s[v]
	return ok
}

func (s Set) Equals(otherF dfa.Fact) bool {
	other := otherF.(Set)
	// am <= other
	for k1 := range s {
		_, ok := other[k1]
		if !ok {
			return false
		}
	}
	// other <= am
	for k1 := range other {
		_, ok := s[k1]
		if !ok {
			return false
		}
	}
	return true
}

func (s Set) String() string {
	variables := make([]string, 0, len(s))
	for k := range s {
		variables = append(variables, k)
	}

	sort.Strings(variables)

	return "{ " + strings.Join(variables, ", ") + " }"
}

func SetFrom(els ...string) Set {
	s := make(Set, len(els))

	for _, el := range els {
		s[el] = struct{}{}
	}

	return s
}

func EmptySet() Set {
	return SetFrom()
}

func PrettyPrintFlows(graph *cfg.CFG, in, out map[int]Set) string {
	res := ""

	ids, idToNode := NodePIFromCFG(graph)

	for _, id := range ids {
		node := idToNode[id]
		res += "\n"
		res += in[node.Label()].String() + "\n"
		res += node.(*NodePI).CFGNode().String() + "\n"
		res += out[node.Label()].String() + "\n"
	}

	return res
}

package ast

import (
	"fmt"
	"github.com/skius/stringlang/internal/frontend/token"
	"sort"
	"strings"
)

func CheckSize(m map[Var]Val) (total int64) {
	for k, v := range m {
		total += int64(len(k)) + int64(len(v))
	}
	return
}

func BoolOf(v Val) bool {
	return v != "false" && v != ""
}

func HasSideEffects(e Expr) bool {
	switch val := e.(type) {
	case Program:
		panic("Program has side effects?")
	case Block:
		any := false
		for _, e := range val {
			if HasSideEffects(e) {
				any = true
				break
			}
		}
		return any
	case Assn:
		return true
	case Var:
		return false
	case Val:
		return false
	case BinOp:
		return HasSideEffects(val.Lhs) || HasSideEffects(val.Rhs)
	case While:
		return HasSideEffects(val.Cond) || HasSideEffects(val.Body)
	case IfElse:
		return HasSideEffects(val.Cond) || HasSideEffects(val.Then) || HasSideEffects(val.Else)
	case Call:
		return true
	case Index:
		return HasSideEffects(val.Source) || HasSideEffects(val.I)
	}
	return false
}

func DefinedVars(e Expr) Set {
	defs := make(Set)
	setDefs(e, defs)
	return defs
}

func setDefs(expr Expr, defs Set) {
	switch val := expr.(type) {
	case Program:
		setDefs(val.Code, defs)
	case Block:
		for _, e := range val {
			setDefs(e, defs)
		}
	case Assn:
		defs[string(val.V)] = struct{}{}
		setDefs(val.E, defs)
	case Var:
		return
	case Val:
		return
	case Arg:
		return
	case BinOp:
		setDefs(val.Lhs, defs)
		setDefs(val.Rhs, defs)
	case While:
		setDefs(val.Cond, defs)
		setDefs(val.Body, defs)
	case IfElse:
		setDefs(val.Cond, defs)
		setDefs(val.Then, defs)
		setDefs(val.Else, defs)
	case Call:
		for _, e := range val.Args {
			setDefs(e, defs)
		}
		// Because we allow arbitrary sources for a call, we need to take those into account
		setDefs(val.Fn, defs)
	case Index:
		setDefs(val.Source, defs)
		setDefs(val.I, defs)
	case Lambda:
		// A lambda defines no variables for its parent scope
		return
	}
	return
}

func UsedBeforeDefVars(e Expr, funcNames Set) Set {
	used := make(Set)
	setUsedBeforeDef(e, used, funcNames)
	return used
}

func setUsedBeforeDef(expr Expr, used Set, funcNames Set) {
	switch val := expr.(type) {
	case Program:
		newFuncNames := make(Set, len(val.Funcs))
		for _, f := range val.Funcs {
			newFuncNames.Union(SetFrom(f.Identifier))
		}
		setUsedBeforeDef(val.Code, used, newFuncNames.Union(funcNames))
	case Block:
		defs := make(Set)
		tempUsed := used.Copy()
		for _, e := range val {
			setUsedBeforeDef(e, tempUsed, funcNames)
			used.Union(tempUsed.Except(defs))
			defs.Union(DefinedVars(e))
		}
	case Assn:
		setUsedBeforeDef(val.E, used, funcNames)
	case Var:
		used[string(val)] = struct{}{}
	case Val:
		return
	case Arg:
		return
	case BinOp:
		setUsedBeforeDef(val.Lhs, used, funcNames)
		setUsedBeforeDef(val.Rhs, used, funcNames)
		// Lhs may have defined some variables for Rhs to use,
		// but let's not, this way we can keep evaluation order ambiguous (if we want to)
	case While:
		setUsedBeforeDef(val.Cond, used, funcNames)
		bodyUsed := UsedBeforeDefVars(val.Body, funcNames)
		// Cond may have defined some variables for Body to use
		used.Union(bodyUsed.Except(DefinedVars(val.Cond)))
	case IfElse:
		setUsedBeforeDef(val.Cond, used, funcNames)
		thenUsed := used.Copy()
		elseUsed := used.Copy()
		setUsedBeforeDef(val.Then, thenUsed, funcNames)
		setUsedBeforeDef(val.Else, elseUsed, funcNames)
		// Cond may have defined some variables for Body to use
		condDefs := DefinedVars(val.Cond)
		used.Union(thenUsed.Union(elseUsed)).Except(condDefs)
	case Call:
		for _, e := range val.Args {
			setUsedBeforeDef(e, used, funcNames)
		}
		// Because we allow arbitrary sources for a call, we need to take those into account
		if variable, ok := val.Fn.(Var); ok {
			varStr := string(variable)
			if funcNames.Contains(varStr) {
				// Calling a named function, hence var is not a usedBeforeDef variable
			} else {
				used[varStr] = struct{}{}
			}
		} else {
			setUsedBeforeDef(val.Fn, used, funcNames)
		}
	case Index:
		setUsedBeforeDef(val.Source, used, funcNames)
		setUsedBeforeDef(val.I, used, funcNames)
	case Lambda:
		innerUsed := UsedBeforeDefVars(val.Code, funcNames)
		// Lambda's used vars are "used \union (innerUsed \except params)
		innerUsed.Except(SetFrom(val.Params...))
		used.Union(innerUsed)
	}
	return
}

func FreeVars(e Expr) map[string]struct{} {
	// I haven't decided yet what the difference should be to UsedVars
	used := make(map[string]struct{})
	setUsedVars(e, used)
	return used
}

func UsedVars(b Block) map[string]struct{} {
	used := make(map[string]struct{})
	setUsedVars(b, used)
	return used
}

func setUsedVars(expr Expr, used map[string]struct{}) {
	switch val := expr.(type) {
	case Program:
		setUsedVars(val.Code, used)
	case Block:
		for _, e := range val {
			setUsedVars(e, used)
		}
	case Assn:
		// used[string(val.V)] = struct{}{} -- the assigned Var isn't "used"
		setUsedVars(val.E, used)
	case Var:
		used[string(val)] = struct{}{}
	case Val:
		return
	case Arg:
		return
	case BinOp:
		setUsedVars(val.Lhs, used)
		setUsedVars(val.Rhs, used)
	case While:
		setUsedVars(val.Cond, used)
		setUsedVars(val.Body, used)
	case IfElse:
		setUsedVars(val.Cond, used)
		setUsedVars(val.Then, used)
		setUsedVars(val.Else, used)
	case Call:
		for _, e := range val.Args {
			setUsedVars(e, used)
		}
		// Because we allow arbitrary sources for a call, we need to take those into account
		setUsedVars(val.Fn, used)
	case Index:
		setUsedVars(val.Source, used)
		setUsedVars(val.I, used)
	case Lambda:
		innerUsed := UsedVars(val.Code)
		// Lambda's used vars are "used \union (innerUsed \except params)
		for _, v := range val.Params {
			delete(innerUsed, v)
		}
		for v := range innerUsed {
			used[v] = struct{}{}
		}
	}
	return
}

func attribToString(a Attrib) string {
	return string(a.(*token.Token).Lit)
}

const (
	SigExternalExit = iota + 1
	SigOutOfMemory
)

// checkExit returns true if we need to exit
func checkExit(c *Context) bool {
	if c.limitStackSize && CheckSize(c.VariableMap) > c.MaxStackSize {
		fmt.Println("Ran out of stack space!")
		select {
		case c.exitChannel <- SigOutOfMemory:
		default:
		}
		return true
	}
	select {
	case <-c.exitChannel:
		// c.exitChannel <- sig // No need to propagate I think?
		return true
	default:
		return false
	}
}

/*
Complete switch of expressions

	switch val := expr.(type) {
	case Program:
	case Block:
	case Assn:
	case Var:
	case Val:
	case Arg:
	case Index:
	case BinOp:
	case IfElse:
	case While:
	case Call:
	case Lambda:
	}

*/



type Set map[string]struct{}

func (s Set) Add(els ...string) {
	for _, el := range els {
		s[el] = struct{}{}
	}
}

func (s Set) UnionCopy(s2 Set) Set {
	return s.Copy().Union(s2)
}

func (s Set) Union(s2 Set) Set {
	for k, v := range s2 {
		s[k] = v
	}

	return s
}

func (s Set) IntersectCopy(s2 Set) Set {
	return s.Copy().Except(s2)
}

func (s Set) Intersect(s2 Set) Set {
	for k := range s {
		if _, ok := s2[k]; !ok {
			delete(s, k)
		}
	}

	return s
}

func (s Set) ExceptCopy(s2 Set) Set {
	return s.Copy().Except(s2)
}

func (s Set) Except(s2 Set) Set {
	for k := range s2 {
		delete(s, k)
	}

	return s
}

func (s Set) Contains(v string) bool {
	_, ok := s[v]
	return ok
}

func (s Set) Equals(s2 Set) bool {
	// am <= s2
	for k1 := range s {
		_, ok := s2[k1]
		if !ok {
			return false
		}
	}
	// s2 <= am
	for k1 := range s2 {
		_, ok := s[k1]
		if !ok {
			return false
		}
	}
	return true
}

func (s Set) Copy() Set {
	s2 := make(Set, len(s))
	for k := range s {
		s2[k] = struct{}{}
	}
	return s2
}

func (s Set) String() string {
	els := make([]string, 0, len(s))
	for k := range s {
		els = append(els, k)
	}

	sort.Strings(els)

	return "{ " + strings.Join(els, ", ") + " }"
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

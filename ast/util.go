package ast

import (
	"github.com/skius/stringlang/token"
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
	case And:
		return HasSideEffects(val.A) || HasSideEffects(val.B)
	case Or:
		return HasSideEffects(val.A) || HasSideEffects(val.B)
	case NotEquals:
		return HasSideEffects(val.A) || HasSideEffects(val.B)
	case Equals:
		return HasSideEffects(val.A) || HasSideEffects(val.B)
	case Concat:
		return HasSideEffects(val.A) || HasSideEffects(val.B)
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
	case And:
		setUsedVars(val.A, used)
		setUsedVars(val.B, used)
	case Or:
		setUsedVars(val.A, used)
		setUsedVars(val.B, used)
	case NotEquals:
		setUsedVars(val.A, used)
		setUsedVars(val.B, used)
	case Equals:
		setUsedVars(val.A, used)
		setUsedVars(val.B, used)
	case Concat:
		setUsedVars(val.A, used)
		setUsedVars(val.B, used)
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
	case Index:
		setUsedVars(val.Source, used)
		setUsedVars(val.I, used)
	}
	return
}

func attribToString(a Attrib) string {
	return string(a.(*token.Token).Lit)
}
func unescape(s string) string {
	in := []rune(s)
	out := make([]rune, 0, len(in))
	var escape bool
	for _, r := range in {
		switch {
		case escape:
			switch r {
			case 'n':
				out = append(out, '\n')
			default:
				out = append(out, r)
			}
			escape = false
		case r == '\\':
			escape = true
		default:
			out = append(out, r)
		}
	}
	return string(out)
}

const (
	SigExternalExit = iota + 1
	SigOutOfMemory
)

// checkExit returns true if we need to exit
func checkExit(c *Context) bool {
	if c.limitStackSize && CheckSize(c.VariableMap) > c.MaxStackSize {
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
	case Concat:
	case And:
	case Or:
	case Equals:
	case NotEquals:
	case IfElse:
	case While:
	case Call:
	}

*/

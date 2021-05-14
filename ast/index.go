package ast

import "strconv"

type Index struct {
	Source Expr
	I      Expr
}

func NewIndex(s, i Attrib) (Expr, error) {
	return Index{Source: s.(Expr), I: i.(Expr)}, nil
}
func NewIndexInt(s, i Attrib) (Expr, error) {
	return Index{Source: s.(Expr), I: Val(attribToString(i))}, nil
}
func (i Index) Eval(c *Context) Val {
	src := string(i.Source.Eval(c))
	idx, err := strconv.Atoi(string(i.I.Eval(c)))
	if err != nil {
		return Val("")
	}
	if idx >= len(src) {
		return Val("")
	}
	return Val(src[idx])
}
func (i Index) String() string {
	srcStr := i.Source.String()
	if i.Source.Precedence() < LeafPrecedence {
		srcStr = "(" + srcStr + ")"
	}
	return srcStr + "[" + i.I.String() + "]"
}
func (i Index) Precedence() int {
	// Leaf, not operator
	return LeafPrecedence
}

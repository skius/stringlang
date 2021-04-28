package ast

import "strings"

type IfElse struct {
	Cond Expr
	Then Expr
	Else Expr
}

func NewIfElse(c, t, e Attrib) (Expr, error) {
	co := c.(Expr)
	th := t.(Expr)
	el := e.(Expr)
	return IfElse{Cond: co, Then: th, Else: el}, nil
}
func (e IfElse) Eval(c *Context) Val {
	if BoolOf(e.Cond.Eval(c)) {
		return e.Then.Eval(c)
	} else {
		return e.Else.Eval(c)
	}
}
func (e IfElse) String() string {
	thenLines := strings.Split(e.Then.String(), "\n")
	thenStr := strings.Join(thenLines, "\n\t")

	elseLines := strings.Split(e.Else.String(), "\n")
	elseStr := strings.Join(elseLines, "\n\t")

	str := "if (" + e.Cond.String() + ") {\n\t" + thenStr + "\n} else {\n\t" + elseStr + "\n}"
	return str
}

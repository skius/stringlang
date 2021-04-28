package ast

import "strings"

type While struct {
	Cond Expr
	Body Expr
}

func NewWhile(c, b Attrib) (Expr, error) {
	co := c.(Expr)
	bo := b.(Expr)
	return While{Cond: co, Body: bo}, nil
}
func (e While) Eval(c *Context) Val {
	var cond Val = e.Cond.Eval(c)
	var body Val
	steps := 0
	for BoolOf(cond) {
		body = e.Body.Eval(c)
		cond = e.Cond.Eval(c)

		if checkExit(c) {
			break
		}
		steps++
	}
	return body
}
func (e While) String() string {
	thenLines := strings.Split(e.Body.String(), "\n")
	thenStr := strings.Join(thenLines, "\n\t")

	str := "while (" + e.Cond.String() + ") {\n\t" + thenStr + "\n}"
	return str
}

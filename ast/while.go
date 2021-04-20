package ast

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
	for boolOf(cond) {
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
	str := "while (" + e.Cond.String() + ") {\n\t" + e.Body.String() + "\n}"
	return str
}

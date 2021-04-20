package ast

type NotEquals struct {
	A Expr
	B Expr
}

func NewNotEquals(a, b Attrib) (Expr, error) {
	return NotEquals{A: a.(Expr), B: b.(Expr)}, nil
}
func (e NotEquals) Eval(c *Context) Val {
	val := "false"
	if e.A.Eval(c) != e.B.Eval(c) {
		val = "true"
	}
	return Val(val)
}
func (e NotEquals) String() string {
	return e.A.String() + " != " + e.B.String()
}

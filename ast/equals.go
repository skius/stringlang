package ast

type Equals struct {
	A Expr
	B Expr
}

func NewEquals(a, b Attrib) (Expr, error) {
	return Equals{A: a.(Expr), B: b.(Expr)}, nil
}
func (e Equals) Eval(c *Context) Val {
	val := "false"
	if e.A.Eval(c) == e.B.Eval(c) {
		val = "true"
	}
	return Val(val)
}
func (e Equals) String() string {
	return e.A.String() + " == " + e.B.String()
}

package ast

type Or struct {
	A Expr
	B Expr
}

func NewOr(a, b Attrib) (Expr, error) {
	return Or{A: a.(Expr), B: b.(Expr)}, nil
}
func (o Or) Eval(c *Context) Val {
	if BoolOf(o.A.Eval(c)) || BoolOf(o.B.Eval(c)) {
		return Val("true")
	} else {
		return Val("false")
	}
}
func (o Or) String() string {
	return o.A.String() + " || " + o.B.String()
}

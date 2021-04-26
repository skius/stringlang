package ast

type And struct {
	A Expr
	B Expr
}

func NewAnd(a, b Attrib) (Expr, error) {
	return And{A: a.(Expr), B: b.(Expr)}, nil
}
func (a And) Eval(c *Context) Val {
	if BoolOf(a.A.Eval(c)) && BoolOf(a.B.Eval(c)) {
		return Val("true")
	} else {
		return Val("false")
	}
}
func (a And) String() string {
	return a.A.String() + " && " + a.B.String()
}

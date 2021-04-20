package ast

type Concat struct {
	A Expr
	B Expr
}

func NewConcat(a, b Attrib) (Expr, error) {
	return Concat{A: a.(Expr), B: b.(Expr)}, nil
}
func (cc Concat) Eval(c *Context) Val {
	return cc.A.Eval(c) + cc.B.Eval(c)
}
func (cc Concat) String() string {
	return cc.A.String() + " + " + cc.B.String()
}

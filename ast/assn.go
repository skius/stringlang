package ast

type Assn struct {
	V Var
	E Expr
}

func NewAssn(v, e Attrib) (Expr, error) {
	va := v.(Var)
	ex := e.(Expr)
	return Assn{V: va, E: ex}, nil
}
func (a Assn) Eval(c *Context) Val {
	newVal := a.E.Eval(c)
	c.VariableMap[a.V] = newVal
	return newVal
}
func (a Assn) String() string {
	return a.V.String() + " = " + a.E.String()
}
func (a Assn) Precedence() int {
	return 0
}

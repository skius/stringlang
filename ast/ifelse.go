package ast

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
	if boolOf(e.Cond.Eval(c)) {
		return e.Then.Eval(c)
	} else {
		return e.Else.Eval(c)
	}
}
func (e IfElse) String() string {
	str := "if (" + e.Cond.String() + ") {\n\t" + e.Then.String() + "\n} else {\n\t" + e.Else.String() + "\n}"
	return str
}

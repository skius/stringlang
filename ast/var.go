package ast

type Var string

func NewVar(a Attrib) (Expr, error) {
	return Var(attribToString(a)), nil
}
// TODO: Make Vars and FuncDecls/Calls be linked: if you call a Var that isn't a function, interpret it as one
// Multiple possibilities: allow reusing same context, maybe instead of fun(a, b, c) args you just use $0, $1, $3 etc
func (v Var) Eval(c *Context) Val {
	return c.VariableMap[v]
}
func (v Var) String() string {
	return string(v)
}

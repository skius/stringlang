package ast

type Val string

func NewVal(a Attrib) (Expr, error) {
	quoted := attribToString(a)
	unquoted := quoted[1 : len(quoted)-1]
	unescaped := unescape(unquoted)
	return Val(unescaped), nil
}
func (v Val) Eval(c *Context) Val {
	return v
}
func (v Val) String() string {
	return "\"" + string(v) + "\""
}

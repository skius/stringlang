package ast

import "strconv"

type Val string

func NewVal(a Attrib) (Expr, error) {
	quoted := attribToString(a)
	res, err := strconv.Unquote(quoted)
	if err != nil {
		panic(err)
	}
	return Val(res), nil
}
func (v Val) Eval(c *Context) Val {
	return v
}
func (v Val) String() string {
	return strconv.Quote(string(v))
}
func (v Val) Precedence() int {
	// Leaf, not operator
	return LeafPrecedence
}

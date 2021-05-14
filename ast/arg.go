package ast

import "strconv"

type Arg int

func NewArg(i Attrib) (Expr, error) {
	s := attribToString(i)
	intValue, err := strconv.Atoi(s)
	return Arg(intValue), err
}
func (a Arg) Eval(c *Context) Val {
	if int(a) >= len(c.Args) {
		return ""
	}
	return Val(c.Args[a])
}
func (a Arg) String() string {
	return "$" + strconv.Itoa(int(a))
}
func (a Arg) Precedence() int {
	// Leaf, not operator
	return LeafPrecedence
}

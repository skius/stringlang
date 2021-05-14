package ast

import (
	"sort"
	"strings"
)

type Lambda struct {
	Params []string
	Code   Block
}

func NewLambda(ps, b Attrib) (Expr, error) {
	params := ps.([]string)
	code := b.(Block)
	return Lambda{Params: params, Code: code}, nil
}

func (l Lambda) Eval(c *Context) Val {
	// Closure by value, copy the value of l's body's var that are used before defined in the current context into its body
	fvMap := UsedBeforeDefVars(l, c.FuncNames())

	// Need to sort as slice because key traversal in maps is non-deterministic
	fv := make([]string, 0, len(fvMap))
	for v := range fvMap {
		fv = append(fv, v)
	}
	sort.Strings(fv)

	captures := make([]Expr, len(fv))
	for i := range fv {
		captures[i] = Assn{V: Var(fv[i]), E: c.VariableMap[Var(fv[i])]}
	}
	l.Code = append(captures, l.Code...)
	// We are evaluating the lambda itself, not calling it, hence we must return the string-value of a lambda
	// For calling, see Lambda.Call
	return Val(l.String())
}

func (l Lambda) String() string {
	res := "fun(" + strings.Join(l.Params, ", ") + ") {\n\t"
	codeLines := strings.SplitN(l.Code.String(), "\n", -1)
	res += strings.Join(codeLines, "\n\t")
	res += "\n}"
	//res = strings.ReplaceAll(res, "\"", "\\\"")
	return res
}

func (l Lambda) Precedence() int {
	return LeafPrecedence
}

func (l Lambda) Call(c *Context, args []Val) Val {
	// Construct corresponding FuncDecl and call that instead
	fDecl := FuncDecl{
		Params:     l.Params,
		Code:       l.Code,
		Identifier: "temp_lambda",
	}
	return fDecl.Call(c, args)
}

package ast

import "strings"

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
	// Closure by value, copy the value of l's body's free vars in the current context into its body
	fv := FreeVars(l)
	for v := range fv {
		l.Code = append([]Expr{Assn{V: Var(v), E: c.VariableMap[Var(v)]}}, l.Code...)
	}
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

func (l Lambda) Call(c *Context, args []Val) Val {
	// Construct corresponding FuncDecl and call that instead
	fDecl := FuncDecl{
		Params:     l.Params,
		Code:       l.Code,
		Identifier: "temp_lambda",
	}
	return fDecl.Call(c, args)
}

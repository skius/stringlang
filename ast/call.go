package ast

import "strings"

type Call struct {
	Fn   Var
	Args CallArgs
}

func NewCall(f, as Attrib) (Expr, error) {
	fn := f.(Var)
	args := as.(CallArgs)
	return Call{Fn: fn, Args: args}, nil
}
func (ca Call) Eval(c *Context) Val {
	if checkExit(c) {
		return ""
	}

	userFn, ok := c.UserFunctionMap[string(ca.Fn)]
	if !ok {
		fn, ok := c.FunctionMap[string(ca.Fn)]
		if !ok {
			panic("function '" + string(ca.Fn) + "' not found.")
		}
		vals := make([]string, 0, len(ca.Args))
		for _, argExp := range ca.Args {
			v := argExp.Eval(c)
			vals = append(vals, string(v))
		}
		res := fn(vals)
		return Val(res)
	}

	vals := make([]Val, 0, len(ca.Args))
	for _, argExp := range ca.Args {
		v := argExp.Eval(c)
		vals = append(vals, v)
	}
	res := userFn.Call(c, vals)
	return res
}
func (ca Call) String() string {
	args := make([]string, 0, len(ca.Args))
	for _, arg := range ca.Args {
		args = append(args, arg.String())
	}

	return ca.Fn.String() + "(" + strings.Join(args, ", ") + ")"
}


// CallArgs is not an Expr, since it can never appear on its own
type CallArgs []Expr

func NewCallArgs() (CallArgs, error) {
	return []Expr{}, nil
}
func CallArgsPrepend(e, a Attrib) (CallArgs, error) {
	args := a.(CallArgs)
	exp := e.(Expr)
	args2 := append([]Expr{exp}, args...)
	return args2, nil
}

// FuncDecl is not an Expr, since it can never appear on its own
type FuncDecl struct {
	Params     []string
	Code       Block
	Identifier string
}

func NewFuncDecl(i, p, b Attrib) (FuncDecl, error) {
	id := attribToString(i)
	params := p.([]string)
	code := b.(Block)
	return FuncDecl{Params: params, Code: code, Identifier: id}, nil
}

const GoStackframeEstimate = 8 * 1024

func (f FuncDecl) Call(c *Context, args []Val) Val {
	newVars := make(map[Var]Val)
	for i, p := range f.Params {
		var argVal Val
		if i < len(args) {
			argVal = args[i]
		}
		newVars[Var(p)] = argVal
	}
	cNew := Context{
		VariableMap:     newVars,
		UserFunctionMap: c.UserFunctionMap,
		FunctionMap:     c.FunctionMap,
		Args:            c.Args,
		MaxStackSize:    c.MaxStackSize - CheckSize(c.VariableMap) - GoStackframeEstimate, // New context needs to account for Go stackframes
		limitStackSize:  c.limitStackSize,
		exitChannel:     c.exitChannel,
	}
	return f.Code.Eval(&cNew)
}
func (f FuncDecl) String() string {
	var id = f.Identifier
	var args = strings.Join(f.Params, ", ")
	codeLines := strings.Split(f.Code.String(), "\n")
	codeStr := strings.Join(codeLines, "\n\t")
	return "fun " + id + "(" + args + ") {\n\t" + codeStr + "\n}"
}
func FuncDeclsAppend(f, fs Attrib) ([]FuncDecl, error) {
	fdecl := f.(FuncDecl)
	funcs := fs.([]FuncDecl)
	return append(funcs, fdecl), nil
}
func FuncParamsPrepend(p, ps Attrib) ([]string, error) {
	param := attribToString(p)
	params := ps.([]string)
	return append([]string{param}, params...), nil
}


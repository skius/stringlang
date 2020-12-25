package ast

import (
	"github.com/skius/stringlang/token"
	"strconv"
	"strings"
)

func NewContext(args []string, funcs map[string]func([]string)string) Context {
	return Context {
		Args: args,
		VariableMap: make(map[Var]Val),
		FunctionMap: funcs,
	}
}

type Context struct {
	Args []string
	VariableMap map[Var]Val
	FunctionMap map[string]func([]string)string
	MaxStackSize int64
	MaxWhileIter int
}

type Attrib interface {}

type Expr interface {
	Eval(Context) Val
	String() string
}

type Val string
func NewVal(a Attrib) (Expr, error) {
	quoted := attribToString(a)
	return Val(quoted[1:len(quoted)-1]), nil
}
func (v Val) Eval(c Context) Val {
	return v
}
func (v Val) String() string {
	return "\"" + string(v) + "\""
}

type Block []Expr
func NewBlock() (Expr, error) {
	return Block([]Expr{}), nil
}
func BlockPrepend(e, b Attrib) (Expr, error) {
	block := b.(Block)
	exp := e.(Expr)
	block2 := append([]Expr{exp}, block...)
	return Block(block2), nil
}
func (b Block) Eval(c Context) Val {
	var last Val
	for _, exp := range b {
		last = exp.Eval(c)
	}
	return last
}
func (b Block) String() string {
	str := ""
	for i, exp := range b {
		if i < len(b) - 1 {
			str += exp.String() + ";\n"
		} else {
			str += exp.String()
		}
	}
	return str
}

type Arg int
func NewArg(i Attrib) (Expr, error) {
	s := attribToString(i)
	intValue, err := strconv.Atoi(s)
	return Arg(intValue), err
}
func (a Arg) Eval(c Context) Val {
	if int(a) > len(c.Args) {
		return ""
	}
	return Val(c.Args[a])
}
func (a Arg) String() string {
	return "$" + strconv.Itoa(int(a))
}

type Equals struct {
	A Expr
	B Expr
}
func NewEquals(a, b Attrib) (Expr, error) {
	return Equals{A: a.(Expr), B: b.(Expr)}, nil
}
func (e Equals) Eval(c Context) Val {
	val := "false"
	if e.A.Eval(c) == e.B.Eval(c) {
		val = "true"
	}
	return Val(val)
}
func (e Equals) String() string {
	return e.A.String() + " == " + e.B.String()
}

type NotEquals struct {
	A Expr
	B Expr
}
func NewNotEquals(a, b Attrib) (Expr, error) {
	return NotEquals{A: a.(Expr), B: b.(Expr)}, nil
}
func (e NotEquals) Eval(c Context) Val {
	val := "false"
	if e.A.Eval(c) != e.B.Eval(c) {
		val = "true"
	}
	return Val(val)
}
func (e NotEquals) String() string {
	return e.A.String() + " != " + e.B.String()
}

type Or struct {
	A Expr
	B Expr
}
func NewOr(a, b Attrib) (Expr, error) {
	return Or{A: a.(Expr), B: b.(Expr)}, nil
}
func (o Or) Eval(c Context) Val {
	if boolOf(o.A.Eval(c)) || boolOf(o.B.Eval(c)) {
		return Val("true")
	} else {
		return Val("false")
	}
}
func (o Or) String() string {
	return o.A.String() + " || " + o.B.String()
}

type And struct {
	A Expr
	B Expr
}
func NewAnd(a, b Attrib) (Expr, error) {
	return And{A: a.(Expr), B: b.(Expr)}, nil
}
func (a And) Eval(c Context) Val {
	if boolOf(a.A.Eval(c)) && boolOf(a.B.Eval(c)) {
		return Val("true")
	} else {
		return Val("false")
	}
}
func (a And) String() string {
	return a.A.String() + " && " + a.B.String()
}

type Concat struct {
	A Expr
	B Expr
}
func NewConcat(a, b Attrib) (Expr, error) {
	return Concat{A: a.(Expr), B: b.(Expr)}, nil
}
func (cc Concat) Eval(c Context) Val {
	return cc.A.Eval(c) + cc.B.Eval(c)
}
func (cc Concat) String() string {
	return cc.A.String() + " + " + cc.B.String()
}

type Var string
func NewVar(a Attrib) (Expr, error) {
	return Var(attribToString(a)), nil
}
func (v Var) Eval(c Context) Val {
	return c.VariableMap[v]
}
func (v Var) String() string {
	return string(v)
}

type Call struct {
	Fn Var
	Args CallArgs
}
func NewCall(f, as Attrib) (Expr, error) {
	fn := f.(Var)
	args := as.(CallArgs)
	return Call {Fn: fn, Args: args}, nil
}
func (ca Call) Eval(c Context) Val {
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
func (ca Call) String() string {
	args := make([]string, 0, len(ca.Args))
	for _, arg := range ca.Args {
		args = append(args, arg.String())
	}

	return ca.Fn.String() + "(" + strings.Join(args, ", ") + ")"
}

type Index struct {
	Source Expr
	I 	Expr
}
func NewIndex(s, i Attrib) (Expr, error) {
	return Index {Source: s.(Expr), I: i.(Expr)}, nil
}
func NewIndexInt(s, i Attrib) (Expr, error) {
	return Index {Source: s.(Expr), I: Val(attribToString(i))}, nil
}
func (i Index) Eval(c Context) Val {
	src := string(i.Source.Eval(c))
	idx, err := strconv.Atoi(string(i.I.Eval(c)))
	if err != nil {
		return Val("")
	}
	if idx >= len(src) {
		return Val("")
	}
	return Val(src[idx])
}
func (i Index) String() string {
	return i.Source.String() + "[" + i.I.String() + "]"
}

type Assn struct {
	V Var
	E Expr
}
func NewAssn(v, e Attrib) (Expr, error) {
	va := v.(Var)
	ex := e.(Expr)
	return Assn {V: va, E: ex}, nil
}
func (a Assn) Eval(c Context) Val {
	newVal := a.E.Eval(c)
	c.VariableMap[a.V] = newVal
	return newVal
}
func (a Assn) String() string {
	return a.V.String() + " = " + a.E.String()
}

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
func (e IfElse) Eval(c Context) Val {
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

type While struct {
	Cond Expr
	Body Expr
}
func NewWhile(c, b Attrib) (Expr, error) {
	co := c.(Expr)
	bo := b.(Expr)
	return While{Cond: co, Body: bo}, nil
}
func (e While) Eval(c Context) Val {
	var cond Val = e.Cond.Eval(c)
	var body Val
	steps := 0
	for boolOf(cond) {
		body = e.Body.Eval(c)
		cond = e.Cond.Eval(c)

		if c.MaxWhileIter > -1 && steps > c.MaxWhileIter {
			break
		}
		if c.MaxStackSize > -1 && checkSize(c.VariableMap) > c.MaxStackSize {
			break
		}
		steps++
	}
	return body
}
func (e While) String() string {
	str := "while (" + e.Cond.String() + ") {\n\t" + e.Body.String() + "\n}"
	return str
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

func attribToString(a Attrib) string {
	return string(a.(*token.Token).Lit)
}
func boolOf(v Val) bool {
	return v != "false" && v != ""
}
func checkSize(m map[Var]Val) (total int64) {
	for k, v := range m {
		total += int64(len(k)) + int64(len(v))
	}
	return
}

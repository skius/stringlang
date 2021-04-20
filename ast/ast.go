package ast

import (
	"strings"
)

type Attrib interface{}

type Expr interface {
	Eval(*Context) Val
	String() string
}

type Program struct {
	Funcs []FuncDecl
	Code  Block
}

func NewProgram(f, b Attrib) (Expr, error) {
	funcs := f.([]FuncDecl)
	code := b.(Block)
	return Program{Funcs: funcs, Code: code}, nil
}

// Useful for building ASTs manually
func EmptyProgram() Program {
	return Program{Funcs: []FuncDecl{}, Code: []Expr{}}
}
func (p Program) Eval(c *Context) Val {
	customFunctions := make(map[string]FuncDecl)
	for _, f := range p.Funcs {
		customFunctions[f.Identifier] = f
	}
	c.UserFunctionMap = customFunctions
	return p.Code.Eval(c)
}
func (p Program) String() string {
	funcdecls := make([]string, len(p.Funcs))
	for i := range p.Funcs {
		funcdecls[i] = p.Funcs[i].String()
	}
	result := strings.Join(funcdecls, "\n")
	return result + "\n" + p.Code.String()
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
func (b Block) Eval(c *Context) Val {
	var last Val
	for _, exp := range b {
		last = exp.Eval(c)
	}
	return last
}
func (b Block) String() string {
	str := ""
	for i, exp := range b {
		if i < len(b)-1 {
			str += exp.String() + ";\n"
		} else {
			str += exp.String()
		}
	}
	return str
}

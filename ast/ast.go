package ast

import (
	"strings"
)

type Attrib interface{}

type Expr interface {
	Eval(*Context) Val
	String() string
	Precedence() int
	//IsRightAssociative() bool // Unnecessary, we only have left-associative expressions currently.
}

const LeafPrecedence int = 100

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
	for _, f := range p.Funcs {
		c.UserFunctionMap[f.Identifier] = f
	}
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
func (p Program) Precedence() int {
	// Unused
	return -1
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
func (b Block) Precedence() int {
	// Unused
	return -1
}

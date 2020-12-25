package stringlang

import (
	"errors"
	"github.com/skius/stringlang/ast"
	"github.com/skius/stringlang/lexer"
	"github.com/skius/stringlang/parser"
)

type Expr = ast.Expr
type Val = ast.Val
type Var = ast.Var
type Context = ast.Context

func NewContext(args []string, funcs map[string]func([]string)string) Context {
	return ast.NewContext(args, funcs)
}

func Parse(body []byte) (ast.Expr, error) {
	l := lexer.NewLexer(body)
	s, err := parser.NewParser().Parse(l)
	if err != nil {
		return nil, err
	}
	e, ok := s.(ast.Expr)
	if !ok {
		return nil, errors.New("couldn't cast parsing result to Expr")
	}
	return e, nil
}

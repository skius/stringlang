package ast

import (
	"fmt"
	"strconv"
)

type Op int

// Op are operators used for BinOp, encoded as their precedence
const (
	OrOp Op = (iota + 1) * 10
	AndOp
	NotEqualsOp
	EqualsOp
	ConcatOp
)

func (o Op) String() string {
	switch o {
	case OrOp:
		return "||"
	case AndOp:
		return "&&"
	case NotEqualsOp:
		return "!="
	case EqualsOp:
		return "=="
	case ConcatOp:
		return "+"
	}
	panic("Op not found:" + strconv.Itoa(int(o)))
}

type BinOp struct {
	Lhs Expr
	Rhs Expr
	Op Op
}

func NewOr(a, b Attrib) (Expr, error) {
	return BinOp{
		Lhs: a.(Expr),
		Rhs: b.(Expr),
		Op:  OrOp,
	}, nil
}
func NewAnd(a, b Attrib) (Expr, error) {
	return BinOp{
		Lhs: a.(Expr),
		Rhs: b.(Expr),
		Op:  AndOp,
	}, nil
}
func NewNotEquals(a, b Attrib) (Expr, error) {
	return BinOp{
		Lhs: a.(Expr),
		Rhs: b.(Expr),
		Op:  NotEqualsOp,
	}, nil
}
func NewEquals(a, b Attrib) (Expr, error) {
	return BinOp{
		Lhs: a.(Expr),
		Rhs: b.(Expr),
		Op:  EqualsOp,
	}, nil
}
func NewConcat(a, b Attrib) (Expr, error) {
	return BinOp{
		Lhs: a.(Expr),
		Rhs: b.(Expr),
		Op:  ConcatOp,
	}, nil
}

func (b BinOp) Eval(c *Context) Val {
	switch b.Op {
	case OrOp:
		return b.evalOr(c)
	case AndOp:
		return b.evalAnd(c)
	case NotEqualsOp:
		return b.evalNotEquals(c)
	case EqualsOp:
		return b.evalEquals(c)
	case ConcatOp:
		return b.evalConcat(c)
	}
	panic("BinOp Op not found")
}
func (b BinOp) String() string {
	lhs := b.Lhs.String()
	rhs := b.Rhs.String()

	// Example:
	// 	Lhs: a == b
	// 	Op: +
	// 	Rhs: c
	// We can't just print "Lhs Op Rhs" <=> "a == b + c", because + has higher precedence, which would "capture" only the b:
	// 	Lhs: a			!! INCORRECT
	// 	Op: ==			!! INCORRECT
	// 	Rhs: b + c      !! INCORRECT
	// Hence we must add parentheses around Lhs

	// For Lhs we need '>', because same precedence operators (i.e. the same operator) can be chained:
	// 	a + b + c <=> (a + b) + c  !! Parentheses unnecessary
	if b.Precedence() > b.Lhs.Precedence() {
		lhs = "(" + lhs + ")"
	}
	// However, because we only have left-associative operators, for the Rhs we must treat the case of same precedence differently,
	// i.e. a + (b + c) may not be printed as a + b + c, because that is equal to (a + b) + c =/= a + (b + c)
	// (Note: If we had different associativity, see here: https://stackoverflow.com/questions/6277747/pretty-print-expression-with-as-few-parentheses-as-possible)
	if b.Precedence() >= b.Rhs.Precedence() {
		rhs = "(" + rhs + ")"
	}
	return fmt.Sprintf("%v %v %v", lhs, b.Op.String(), rhs)
}
func (b BinOp) Precedence() int {
	return int(b.Op)
}


func (b BinOp) evalOr(c *Context) Val {
	if BoolOf(b.Lhs.Eval(c)) || BoolOf(b.Rhs.Eval(c)) {
		return Val("true")
	} else {
		return Val("false")
	}
}
func (b BinOp) evalAnd(c *Context) Val {
	if BoolOf(b.Lhs.Eval(c)) && BoolOf(b.Rhs.Eval(c)) {
		return Val("true")
	} else {
		return Val("false")
	}
}
func (b BinOp) evalNotEquals(c *Context) Val {
	val := "false"
	if b.Lhs.Eval(c) != b.Rhs.Eval(c) {
		val = "true"
	}
	return Val(val)
}
func (b BinOp) evalEquals(c *Context) Val {
	val := "false"
	if b.Lhs.Eval(c) == b.Rhs.Eval(c) {
		val = "true"
	}
	return Val(val)
}
func (b BinOp) evalConcat(c *Context) Val {
	return b.Lhs.Eval(c) + b.Rhs.Eval(c)
}

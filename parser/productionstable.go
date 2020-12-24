// Code generated by gocc; DO NOT EDIT.

package parser

import "github.com/skius/stringlang/ast"

type (
	//TODO: change type and variable names to be consistent with other tables
	ProdTab      [numProductions]ProdTabEntry
	ProdTabEntry struct {
		String     string
		Id         string
		NTType     int
		Index      int
		NumSymbols int
		ReduceFunc func([]Attrib) (Attrib, error)
	}
	Attrib interface {
	}
)

var productionsTable = ProdTab{
	ProdTabEntry{
		String: `S' : ExprSeq	<<  >>`,
		Id:         "S'",
		NTType:     0,
		Index:      0,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `ExprSeq : ConcatExpr Expr	<< ast.ExprSeqAppend(X[0], X[1]) >>`,
		Id:         "ExprSeq",
		NTType:     1,
		Index:      1,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.ExprSeqAppend(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `ConcatExpr : ExprSeq ";"	<< X[0], nil >>`,
		Id:         "ConcatExpr",
		NTType:     2,
		Index:      2,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `ConcatExpr : empty	<< ast.NewExprSeq() >>`,
		Id:         "ConcatExpr",
		NTType:     2,
		Index:      3,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewExprSeq()
		},
	},
	ProdTabEntry{
		String: `Expr : Var "=" Expr	<< ast.NewAssn(X[0], X[2]) >>`,
		Id:         "Expr",
		NTType:     3,
		Index:      4,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewAssn(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr : Expr000	<< X[0], nil >>`,
		Id:         "Expr",
		NTType:     3,
		Index:      5,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr000 : Expr000 "||" Expr00	<< ast.NewOr(X[0], X[2]) >>`,
		Id:         "Expr000",
		NTType:     4,
		Index:      6,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewOr(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr000 : Expr00	<< X[0], nil >>`,
		Id:         "Expr000",
		NTType:     4,
		Index:      7,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr00 : Expr00 "&&" Expr01	<< ast.NewAnd(X[0], X[2]) >>`,
		Id:         "Expr00",
		NTType:     5,
		Index:      8,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewAnd(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr00 : Expr01	<< X[0], nil >>`,
		Id:         "Expr00",
		NTType:     5,
		Index:      9,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr01 : Expr01 "!=" Expr0	<< ast.NewNotEquals(X[0], X[2]) >>`,
		Id:         "Expr01",
		NTType:     6,
		Index:      10,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewNotEquals(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr01 : Expr0	<< X[0], nil >>`,
		Id:         "Expr01",
		NTType:     6,
		Index:      11,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr0 : Expr0 "==" Expr1	<< ast.NewEquals(X[0], X[2]) >>`,
		Id:         "Expr0",
		NTType:     7,
		Index:      12,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewEquals(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr0 : Expr1	<< X[0], nil >>`,
		Id:         "Expr0",
		NTType:     7,
		Index:      13,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr1 : Expr1 "+" Expr2	<< ast.NewConcat(X[0], X[2]) >>`,
		Id:         "Expr1",
		NTType:     8,
		Index:      14,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewConcat(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr1 : Expr2	<< X[0], nil >>`,
		Id:         "Expr1",
		NTType:     8,
		Index:      15,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr2 : IfElse	<< X[0], nil >>`,
		Id:         "Expr2",
		NTType:     9,
		Index:      16,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr2 : While	<< X[0], nil >>`,
		Id:         "Expr2",
		NTType:     9,
		Index:      17,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr2 : string_lit	<< ast.NewVal(X[0]) >>`,
		Id:         "Expr2",
		NTType:     9,
		Index:      18,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewVal(X[0])
		},
	},
	ProdTabEntry{
		String: `Expr2 : Arg	<< X[0], nil >>`,
		Id:         "Expr2",
		NTType:     9,
		Index:      19,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr2 : Var	<< X[0], nil >>`,
		Id:         "Expr2",
		NTType:     9,
		Index:      20,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr2 : Var "(" ExprSeqCall ")"	<< ast.NewCall(X[0], X[2]) >>`,
		Id:         "Expr2",
		NTType:     9,
		Index:      21,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewCall(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr2 : "(" Expr ")"	<< X[1], nil >>`,
		Id:         "Expr2",
		NTType:     9,
		Index:      22,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[1], nil
		},
	},
	ProdTabEntry{
		String: `Expr2 : Index	<< X[0], nil >>`,
		Id:         "Expr2",
		NTType:     9,
		Index:      23,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Index : Expr2 "[" Expr "]"	<< ast.NewIndex(X[0], X[2]) >>`,
		Id:         "Index",
		NTType:     10,
		Index:      24,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewIndex(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Index : Expr2 "[" int_lit "]"	<< ast.NewIndexInt(X[0], X[2]) >>`,
		Id:         "Index",
		NTType:     10,
		Index:      25,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewIndexInt(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `ExprSeqCall : ConcatExprCall Expr	<< ast.ExprSeqAppend(X[0], X[1]) >>`,
		Id:         "ExprSeqCall",
		NTType:     11,
		Index:      26,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.ExprSeqAppend(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `ExprSeqCall : empty	<< ast.NewExprSeq() >>`,
		Id:         "ExprSeqCall",
		NTType:     11,
		Index:      27,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewExprSeq()
		},
	},
	ProdTabEntry{
		String: `ConcatExprCall : ExprSeqCall ","	<< X[0], nil >>`,
		Id:         "ConcatExprCall",
		NTType:     12,
		Index:      28,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `ConcatExprCall : empty	<< ast.NewExprSeq() >>`,
		Id:         "ConcatExprCall",
		NTType:     12,
		Index:      29,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewExprSeq()
		},
	},
	ProdTabEntry{
		String: `Arg : "$" int_lit	<< ast.NewArg(X[1]) >>`,
		Id:         "Arg",
		NTType:     13,
		Index:      30,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewArg(X[1])
		},
	},
	ProdTabEntry{
		String: `Var : id	<< ast.NewVar(X[0]) >>`,
		Id:         "Var",
		NTType:     14,
		Index:      31,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewVar(X[0])
		},
	},
	ProdTabEntry{
		String: `IfElse : "if" "(" Expr ")" "{" ExprSeq "}" "else" "{" ExprSeq "}"	<< ast.NewIfElse(X[2], X[5], X[9]) >>`,
		Id:         "IfElse",
		NTType:     15,
		Index:      32,
		NumSymbols: 11,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewIfElse(X[2], X[5], X[9])
		},
	},
	ProdTabEntry{
		String: `IfElse : "if" "(" Expr ")" "{" ExprSeq "}" "else" IfElse	<< ast.NewIfElse(X[2], X[5], X[8]) >>`,
		Id:         "IfElse",
		NTType:     15,
		Index:      33,
		NumSymbols: 9,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewIfElse(X[2], X[5], X[8])
		},
	},
	ProdTabEntry{
		String: `While : "while" "(" Expr ")" "{" ExprSeq "}"	<< ast.NewWhile(X[2], X[5]) >>`,
		Id:         "While",
		NTType:     16,
		Index:      34,
		NumSymbols: 7,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewWhile(X[2], X[5])
		},
	},
}

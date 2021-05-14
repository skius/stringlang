package optimizer

import (
	. "github.com/skius/stringlang/ast"
	"strconv"
)

type normalizer struct {
	names map[string]struct{} // Used variables
	last  int                 // last generated temp
}

// Normalize returns an equivalent program, but with at most one level of expression nesting
// Useful as a first step for further optimizations
func Normalize(prog Program) Program {
	n := new(normalizer)
	n.reset()
	n.names = UsedVars(prog.Code)
	code := n.compileStmt(prog.Code)

	funcs := make([]FuncDecl, len(prog.Funcs))
	for i := range prog.Funcs {
		n.reset()
		n.names = UsedVars(prog.Funcs[i].Code)
		funcs[i] = FuncDecl{
			Params:     prog.Funcs[i].Params,
			Code:       n.compileStmt(prog.Funcs[i].Code),
			Identifier: prog.Funcs[i].Identifier,
		}
	}
	return Program{Funcs: funcs, Code: code}
}

// returns list of expressions that are equivalent to root-level expression e
func (n *normalizer) compileStmt(_e Expr) (res []Expr) {
	res = []Expr{}
	switch e := _e.(type) {
	case Program:
		panic("not impl")
	case Block:
		for _, subExpr := range e {
			res = append(res, n.compileStmt(subExpr)...)
		}
	case Assn:
		code, loc := n.compileExpr(e.E, 1)
		res = code
		res = append(res, Assn{V: e.V, E: loc})
	case Var:
		res = []Expr{e}
	case Val:
		res = []Expr{e}
	case Arg:
		res = []Expr{e}
	case Index:
		codeSrc, locSrc := n.compileExpr(e.Source, 0)
		codeI, locI := n.compileExpr(e.I, 0)
		res = codeSrc
		res = append(res, codeI...)
		res = append(res, Index{Source: locSrc, I: locI})
	case BinOp:
		codeLhs, locLhs := n.compileExpr(e.Lhs, 0)
		codeRhs, locRhs := n.compileExpr(e.Rhs, 0)
		res = codeLhs
		res = append(res, codeRhs...)
		res = append(res, BinOp{Lhs: locLhs, Rhs: locRhs, Op: e.Op})
	case IfElse: // TODO: currently short-circuiting is broken
		codeCond, locCond := n.compileExpr(e.Cond, 0)
		res = append(res, codeCond...)

		newIfElse := IfElse{
			Cond: locCond,
			Then: Block(n.compileStmt(e.Then)),
			Else: Block(n.compileStmt(e.Else)),
		}
		res = append(res, newIfElse)
	case While: // TODO: this is a tricky one, I believe...
		// Need to put excess condition code _before_ loop but also _append to the end_ of the body to keep semantics
		codeCond, locCond := n.compileExpr(e.Cond, 0)
		res = append(res, codeCond...)

		newWhile := While{
			Cond: locCond,
			Body: Block(append(n.compileStmt(e.Body), codeCond...)),
		}
		res = append(res, newWhile)
	case Call:
		var newArgs CallArgs
		for _, argE := range e.Args {
			code, loc := n.compileExpr(argE, 0)
			res = append(res, code...)
			newArgs = append(newArgs, loc)
		}
		res = append(res, Call{Fn: e.Fn, Args: newArgs})
	}
	return
}

// returns expression that is equivalent to expression e when used after evaluating returned list
// maxDepth allows forcing of primitive expressions, e.g. recursively calling compileExpr with maxDepth = 0 for "a + b"
// (a + b)["2"] -> code = [__a = a + b], loc = __a["2"]
func (n *normalizer) compileExpr(_e Expr, maxDepth int) (code []Expr, loc Expr) {
	code = []Expr{}
	switch e := _e.(type) {
	case Program:
		panic("Program not subexpr")
	case Block:
		panic("Block not subexpr")
	case Assn:
		loc = e.V
		code = n.compileStmt(e)
	case Var:
		loc = e
	case Val:
		loc = e
	case Arg:
		loc = e
	case Index:
		if maxDepth == 0 {
			v := Var(n.genName())
			code = n.compileStmt(Assn{V: v, E: e})
			loc = v
		} else if maxDepth == 1 {
			codeSrc, locSrc := n.compileExpr(e.Source, 0)
			codeI, locI := n.compileExpr(e.I, 0)
			code = append(code, codeSrc...)
			code = append(code, codeI...)
			loc = Index{Source: locSrc, I: locI}
		} else {
			panic("maxDepth > 1")
		}
	case BinOp:
		if maxDepth == 0 {
			v := Var(n.genName())
			code = n.compileStmt(Assn{V: v, E: e})
			loc = v
		} else if maxDepth == 1 {
			codeLhs, locLhs := n.compileExpr(e.Lhs, 0)
			codeRhs, locRhs := n.compileExpr(e.Rhs, 0)
			code = append(code, codeLhs...)
			code = append(code, codeRhs...)
			loc = BinOp{Lhs: locLhs, Rhs: locRhs, Op: e.Op}
		} else {
			panic("maxDepth > 1")
		}
	// Always compile control flow to root-level expressions
	case IfElse:
		v := Var(n.genName())

		then := e.Then.(Block)
		newThen := make(Block, len(then))
		copy(newThen, then)
		lastThen := len(then) - 1
		newThen[lastThen] = Assn{V: v, E: newThen[lastThen]}

		eelse := e.Else.(Block)
		newElse := make(Block, len(eelse))
		copy(newElse, eelse)
		lastElse := len(eelse) - 1
		newElse[lastElse] = Assn{V: v, E: newElse[lastElse]}

		newIfElse := IfElse{Cond: e.Cond, Then: newThen, Else: newElse}
		loc = v
		code = n.compileStmt(newIfElse)
	case While:
		v := Var(n.genName())

		body := e.Body.(Block)
		newBody := make(Block, len(body))
		copy(newBody, body)
		lastBody := len(body) - 1
		newBody[lastBody] = Assn{V: v, E: body[lastBody]}

		// Because the v might not be initialized if the while loop never executes, we need to initialize it before the while loop
		init := Assn{V: v, E: Val("")}
		code = []Expr{init}

		loc = v
		code = append(code, n.compileStmt(While{Cond: e.Cond, Body: newBody})...)
	case Call:
		if maxDepth == 0 {
			v := Var(n.genName())
			code = n.compileStmt(Assn{V: v, E: e})
			loc = v
		} else if maxDepth == 1 {
			var newArgs CallArgs
			for _, argE := range e.Args {
				codeArg, locArg := n.compileExpr(argE, 0)
				code = append(code, codeArg...)
				newArgs = append(newArgs, locArg)
			}
			loc = Call{Fn: e.Fn, Args: newArgs}
		} else {
			panic("maxDepth > 1")
		}
	}
	return
}

func (n *normalizer) reset() {
	n.names = make(map[string]struct{})
	n.last = 1
}

func (n *normalizer) genName() string {
	lastNum := strconv.Itoa(n.last)
	_, ok := n.names["__temp"+lastNum]
	for ; ok; _, ok = n.names["__temp"+lastNum] {
		lastNum = strconv.Itoa(n.last)
		n.last++
	}
	n.names["__temp"+lastNum] = struct{}{}
	return "__temp" + lastNum
}

//func succ(s string) string {
//	sRev := reverse(s)
//
//	if sRev == "" {
//		return "a"
//	}
//
//	if sRev[0] == 'z' {
//		return reverse("a" + reverse(succ(reverse(sRev[1:]))))
//	}
//
//	res := string(sRev[0]+1) + sRev[1:]
//	return reverse(res)
//}
//
//func reverse(s string) string {
//	var sRev string
//	for i := len(s) - 1; i >= 0; i-- {
//		sRev += string(s[i])
//	}
//	return sRev
//}

//func max(a, b int) int {
//	if a > b {
//		return a
//	}
//	return b
//}
//
//func primitiveLevel(e Expr) int {
//	switch val := e.(type) {
//	case Program:
//		return primitiveLevel(val.Code)
//	case Block:
//		codeMax := 0
//		for _, e := range val {
//			codeMax = max(primitiveLevel(e), codeMax)
//		}
//		return len(val) + codeMax
//	case Assn:
//		return 1 + primitiveLevel(val.E)
//	case Var:
//		return 0
//	case Val:
//		return 0
//	case BinOp:
//		return 1 + max(primitiveLevel(val.Lhs), primitiveLevel(val.Rhs))
//	case While:
//		return 1 + max(primitiveLevel(val.Cond), primitiveLevel(val.Body))
//	case IfElse:
//		return 1 + max(max(primitiveLevel(val.Cond), primitiveLevel(val.Then)), primitiveLevel(val.Else))
//	case Call:
//		argM := 0
//		for _, e := range val.Args {
//			argM = max(primitiveLevel(e), argM)
//		}
//		return 1 + argM
//	case Index:
//		return 1 + max(primitiveLevel(val.Source), primitiveLevel(val.I))
//	}
//	return 0
//}

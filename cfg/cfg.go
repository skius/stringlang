package cfg

import (
	"fmt"
	"github.com/skius/stringlang/ast"
	"strconv"
)

type Node struct {
	SuccNotTaken  *Node
	SuccTaken     *Node
	PredsNotTaken []*Node
	PredsTaken    []*Node
	FuncSucc      *Node // TODO: Fill these
	Label         int
	Expr          ast.Expr
	Type		  string
}

func (n *Node) IsIf() bool {
	return n.Type == "if"
}

func (n *Node) IsWhile() bool {
	return n.Type == "while"
}

func (n *Node) IsSentinel() bool {
	return n.Expr == nil
}

func (n *Node) IsBranch() bool {
	return n.SuccTaken != nil
}

func (n *Node) String() string {
	if n.IsSentinel() {
		return "<sentinel>"
	}
	lbl := strconv.Itoa(n.Label) + ": "
	if n.Type != "" {
		return lbl + n.Type + " (" + n.Expr.String() + ")"
	}
	return lbl + n.Expr.String()
}

// Remove removes n from the CFG by updating its preds and succs
func (n *Node) Remove() {
	if n.IsSentinel() {
		panic("Trying to remove sentinel from CFG")
	}
	if n.IsWhile() {
		panic("Trying to remove While from CFG")
	}
	if n.SuccNotTaken == nil {
		panic("Trying to remove exit node")
	}

	fmt.Println("Removing, ", *n)



	if n.SuccTaken != nil && n.SuccTaken.Label != n.SuccNotTaken.Label {
		panic("Removing branch with different nottaken/taken branches")
	}

	predsNT := n.PredsNotTaken
	predsT := n.PredsTaken

	succPredsNT := make([]*Node, 0)
	succPredsT := make([]*Node, 0)

	for _, pred := range n.SuccNotTaken.PredsNotTaken {
		if pred.Label != n.Label {
			succPredsNT = append(succPredsNT, pred)
		}
	}
	succPredsNT = append(succPredsNT, predsNT...)

	for _, pred := range n.SuccNotTaken.PredsTaken {
		if pred.Label != n.Label {
			succPredsNT = append(succPredsNT, pred)
		}
	}
	succPredsT = append(succPredsT, predsT...)

	n.SuccNotTaken.PredsNotTaken = predsNT
	n.SuccNotTaken.PredsTaken = predsT

	for _, predNT := range predsNT {
		predNT.SuccNotTaken = n.SuccNotTaken
	}

	for _, predT := range predsT {
		predT.SuccTaken = n.SuccNotTaken
	}






	//predsNT := n.PredsNotTaken
	//predsT := n.PredsTaken
	//
	//// SuccTaken will be removed as well, because it won't be reachable anymore
	//// But we'll do this in a later pass where we compute reachability using just forward edges
	//
	//succ := n.SuccNotTaken
	//
	//succPredsNTNew := make([]*Node, 0, len(succ.PredsNotTaken) + len(predsNT))
	//
	//for _, p := range succ.PredsNotTaken {
	//	if p == n {
	//		continue
	//	}
	//	succPredsNTNew = append(succPredsNTNew, p)
	//}
	//
	//succ.PredsNotTaken = succPredsNTNew
	//
	//for _, predNT := range predsNT {
	//	predNT.SuccNotTaken = succ
	//	succ.PredsNotTaken = append(succ.PredsNotTaken, predNT)
	//}
	//
	//for _, predT := range predsT {
	//	predT.SuccTaken = succ
	//	succ.PredsTaken = append(succ.PredsTaken, predT)
	//}

	// There shouldn't be forward references to n around anymore now
}

type CFG struct {
	Entry *Node
	Exits []*Node
}

type counter struct {
	curr int
}

// NewFromBlock is useful if we do local analysis of the program
func NewFromBlock(b ast.Block) *CFG {
	cfg := new(CFG)
	ctr := new(counter)
	ctr.curr = -1

	sentinel := new(Node)
	sentinel.Label = ctr.incAndGet()
	cfg.Entry = sentinel
	exits := fillBlock(cfg.Entry, b, ctr, false)
	cfg.Exits = exits

	fillPreds(cfg)
	return cfg
}

// New returns the CFG of the top-level expressions and a map of FuncDecls to CFGs
func New(prog ast.Program) (*CFG, map[string]*CFG) {
	cfg := new(CFG)
	ctr := new(counter)
	ctr.curr = -1 // to start our labelling at 0

	sentinel := new(Node)
	sentinel.Label = ctr.incAndGet()
	cfg.Entry = sentinel
	exits := fillBlock(cfg.Entry, prog.Code, ctr, false)

	//sentinelExits := make([]*Node, 0, len(exits))
	//for _, exit := range exits {
	//	sentinelExit := new(Node)
	//	sentinelExit.Label = ctr.incAndGet()
	//	exit.SuccNotTaken = sentinelExit
	//	sentinelExit.PredsNotTaken = []*Node{exit}
	//	sentinelExits = append(sentinelExits, sentinelExit)
	//}
	cfg.Exits = exits

	fillPreds(cfg)

	cfgFuncs := make(map[string]*CFG)

	// Reuse ctr so we have globally unique labels
	// (will cause problems if I implement separate compilation units)
	for _, fd := range prog.Funcs {
		funcCfg := new(CFG)
		sentinel = new(Node)
		sentinel.Label = ctr.incAndGet()
		funcCfg.Entry = sentinel
		funcCfg.Exits = fillBlock(funcCfg.Entry, fd.Code, ctr, false)
		fillPreds(funcCfg)
		cfgFuncs[fd.Identifier] = funcCfg
	}

	return cfg, cfgFuncs
}

func VisitAll(start *Node, f func(*Node)) {
	visited := make(map[int]bool)

	var dfs func(*Node)
	dfs = func(curr *Node) {
		if curr == nil || visited[curr.Label] {
			return
		}
		if curr.IsSentinel() {
			// Skip it, but consider successors
			dfs(curr.SuccNotTaken)
			dfs(curr.SuccTaken)
			return
		}
		visited[curr.Label] = true

		f(curr)

		dfs(curr.SuccNotTaken)
		dfs(curr.SuccTaken)
	}

	dfs(start)
}

// Visit runs the given closure over the CFG in DFS preorder
func (cfg *CFG) Visit(f func(*Node)) {
	VisitAll(cfg.Entry, f)
}

// Fills in backward edges for CFG which already has forward edges
func fillPreds(cfg *CFG) {
	cfg.Visit(func(curr *Node) {
		if curr.SuccNotTaken != nil {
			succ := curr.SuccNotTaken
			succ.PredsNotTaken = append(succ.PredsNotTaken, curr)
		}

		if curr.SuccTaken != nil {
			succ := curr.SuccTaken
			succ.PredsTaken = append(succ.PredsTaken, curr)
		}
	})
}

// Returns exits of block
// Only fills in forward-edges, because backward (pred) edges can be added easily using a visitor
func fillBlock(entryPred *Node, block ast.Block, ctr *counter, isBranch bool) []*Node {
	preds := []*Node{}
	if entryPred != nil {
		preds = []*Node{entryPred}
	}

	updateSucc := func(succ *Node) {
		for _, pred := range preds {
			if isBranch {
				pred.SuccTaken = succ
			} else {
				pred.SuccNotTaken = succ
			}
		}
		isBranch = false
	}

	for _, expr := range block {
		//if exprN, ok := expr.(ExprWithNode); ok {
		//	expr = exprN.Expr
		//}

		switch e := expr.(type) {
		case ast.IfElse:
			condNode := buildNode(e.Cond, ctr)
			condNode.Type = "if"
			updateSucc(condNode)
			tExits := fillBlock(condNode, e.Then.(ast.Block), ctr, true)
			nExits := fillBlock(condNode, e.Else.(ast.Block), ctr, false)
			preds = append(tExits, nExits...)
		case ast.While:
			condNode := buildNode(e.Cond, ctr)
			condNode.Type = "while"
			updateSucc(condNode)
			bExits := fillBlock(condNode, e.Body.(ast.Block), ctr, true)
			for _, pred := range bExits {
				pred.SuccNotTaken = condNode
			}
			preds = []*Node{condNode}
		default:
			n := buildNode(expr, ctr)
			updateSucc(n)
			preds = []*Node{n}
		}
	}

	return preds
}

func buildNode(expr ast.Expr, ctr *counter) *Node {
	n := new(Node)
	n.Label = ctr.incAndGet()
	n.Expr = expr
	return n
}

func (c *counter) incAndGet() int {
	c.curr++
	return c.curr
}

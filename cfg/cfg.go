package cfg

import "github.com/skius/stringlang/ast"

type Node struct {
	SuccNotTaken  *Node
	SuccTaken     *Node
	PredsNotTaken []*Node
	PredsTaken    []*Node
	FuncSucc      *Node // TODO: Fill these
	Label         int
	Expr          ast.Expr
}

type CFG struct {
	Entry *Node
	Exits []*Node
}

type counter struct {
	curr int
}

// New returns the CFG of the top-level expressions and a map of FuncDecls to CFGs
func New(prog ast.Program) (*CFG, map[string]*CFG) {
	cfg := new(CFG)
	ctr := new(counter)

	// Block is non-empty, can do this
	// TODO: ugly
	cfg.Entry = buildNode(prog.Code[0], ctr)
	cfg.Exits = fillBlock(cfg.Entry, prog.Code[1:], ctr, false)

	fillPreds(cfg)

	cfgFuncs := make(map[string]*CFG)

	// Reuse ctr so we have globally unique labels
	// (will cause problems when if I implement separate compilation units)
	for _, fd := range prog.Funcs {
		funcCfg := new(CFG)
		funcCfg.Entry = buildNode(fd.Code[0], ctr)
		funcCfg.Exits = fillBlock(funcCfg.Entry, fd.Code[1:], ctr, false)
		fillPreds(funcCfg)
		cfgFuncs[fd.Identifier] = funcCfg
	}

	return cfg, cfgFuncs
}

// Visit runs the given closure over the CFG in DFS preorder
func (cfg *CFG) Visit(f func(*Node)) {
	visited := make(map[int]bool)

	var dfs func(*Node)
	dfs = func(curr *Node) {
		if curr == nil || visited[curr.Label] {
			return
		}
		visited[curr.Label] = true

		f(curr)

		dfs(curr.SuccNotTaken)
		dfs(curr.SuccTaken)
	}

	dfs(cfg.Entry)
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

		switch e := expr.(type) {
		case ast.IfElse:
			condNode := buildNode(e.Cond, ctr)
			updateSucc(condNode)
			tExits := fillBlock(condNode, e.Then.(ast.Block), ctr, true)
			nExits := fillBlock(condNode, e.Else.(ast.Block), ctr, false)
			preds = append(tExits, nExits...)
		case ast.While:
			condNode := buildNode(e.Cond, ctr)
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

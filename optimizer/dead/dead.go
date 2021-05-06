package dead

import (
	"fmt"
	"github.com/skius/stringlang/ast"
	"github.com/skius/stringlang/cfg"
)

// Eliminate expects a normalized Program and returns the dead-code-eliminated Program
func Eliminate(prog ast.Program) ast.Program {
	block := prog.Code
	funcs := prog.Funcs

	block = eliminateBlock(block)

	for i := range funcs {
		funcs[i].Code = eliminateBlock(funcs[i].Code)
	}

	return ast.Program{
		Code: block,
		Funcs: funcs,
	}
}

func eliminateBlock(b ast.Block) ast.Block {
	g := cfg.NewFromBlock(b)
	// Do three side-effect passes
	changed := true
	for changed {
		prev := g.Collect()
		g = sideeffectDCE(g)
		post := g.Collect()
		changed = post.String() != prev.String()
		if changed {
			fmt.Println()
			fmt.Println("Changed!")
			fmt.Println(prev)
			fmt.Println()
			fmt.Println(post)
			fmt.Println("----")
		}
	}
	//g = sideeffectDCEWhile(g)
	//g = sideeffectDCE(g)
	//g = sideeffectDCE(g)
	//g = sideeffectDCEWhile(g)
	//g = sideeffectDCE(g)
	//g = sideeffectDCE(g)
	//g = sideeffectDCEWhile(g)
	return g.Collect()
}

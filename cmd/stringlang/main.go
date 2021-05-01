package main

import (
	"flag"
	"fmt"
	"github.com/skius/stringlang"
	"github.com/skius/stringlang/ast"
	"github.com/skius/stringlang/cfg"
	"github.com/skius/stringlang/cmd/stringlang/repl"
	"github.com/skius/stringlang/optimizer"
	"github.com/skius/stringlang/optimizer/analysis/liveness"
	"github.com/skius/stringlang/optimizer/analysis/sideeffect"
	"github.com/skius/stringlang/optimizer/analysis/util"
	"io/ioutil"
	"time"
)

func main() {
	var normalize bool
	flag.BoolVar(&normalize, "normalize", false, "Normalize program")

	var printSource bool
	flag.BoolVar(&printSource, "print", false, "Print program source")

	var graphvizFile string
	flag.StringVar(&graphvizFile, "graphviz", "", "Write CFG of program as .dot file to argument")

	var sideffectAnalysis bool
	flag.BoolVar(&sideffectAnalysis, "sideeffect", false, "Print results of side-effect analysis [forces --normalize]")

	var livenessAnalysis bool
	flag.BoolVar(&livenessAnalysis, "liveness", false, "Print results of liveness analysis [forces --normalize]")

	flag.Parse()

	anyFlagSet := false
	flag.Visit(func(f *flag.Flag) {
		anyFlagSet = true
	})

	if !anyFlagSet && len(flag.Args()) == 0 {
		t := repl.DefaultTerminal()
		r := repl.Init(t)
		r.Run()
		return
	}

	if anyFlagSet && len(flag.Args()) == 0 {
		fmt.Println("Usage: ./stringlang [[--normalize] [--graphviz=file] [--print] <program.stringlang> [..<args>]]")
		return
	}

	sourceFile := flag.Args()[0]
	source, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		panic(err)
	}

	expr, err := stringlang.Parse(source)
	if err != nil {
		panic(err)
	}

	program := expr.(ast.Program)

	if sideffectAnalysis {
		// Forces normalize
		normalize = false
		program = optimizer.Normalize(program)

		fmt.Println("Side-effect Analysis:")
		// TODO memoize normalization results, graphs etc using functions
		g, _ := cfg.New(program)
		in, out := sideeffect.Compute(g)
		str := util.PrettyPrintFlows(g, in, out)
		fmt.Println(str)
		fmt.Println()
		fmt.Println()
	}

	if livenessAnalysis {
		// Forces normalize
		normalize = false
		program = optimizer.Normalize(program)

		fmt.Println("Liveness Analysis:")
		g, _ := cfg.New(program)
		in, out := liveness.Compute(g)
		str := util.PrettyPrintFlows(g, in, out)
		fmt.Println(str)
		fmt.Println()
		fmt.Println()
	}

	if normalize {
		program = optimizer.Normalize(program)
	}

	if graphvizFile != "" {
		g, fgs := cfg.New(program)
		graph := cfg.GraphViz(g, fgs)

		output := graph.String()

		err = ioutil.WriteFile(graphvizFile, []byte(output), 0666)
		if err != nil {
			panic(err)
		}
	}

	if printSource {
		fmt.Println(program.String())
		return
	}

	result, err := stringlang.EvalOrTimeout(stringlang.ExampleContext(true), program, time.Second*30)
	if err != nil {
		panic(err)
	}

	fmt.Println("Returns:")
	fmt.Println(result)
}

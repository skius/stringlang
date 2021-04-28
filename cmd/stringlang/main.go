package main

import (
	"flag"
	"fmt"
	"github.com/skius/stringlang"
	"github.com/skius/stringlang/ast"
	"github.com/skius/stringlang/cfg"
	"github.com/skius/stringlang/optimizer"
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

	flag.Parse()

	anyFlagSet := false
	flag.Visit(func(f *flag.Flag) {
		anyFlagSet = true
	})

	if !anyFlagSet && len(flag.Args()) == 0 {
		r := new(Repl)
		r.Run()
		return
	}

	if anyFlagSet && len(flag.Args()) == 0 {
		fmt.Println("Usage: ./stringlang [[--normalize] [--graphviz=file] [--print] <program.stringlang> [..<args>]]")
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

	result, err := evalOrTimeout(exampleContext(true), program, time.Second*30)
	if err != nil {
		panic(err)
	}

	fmt.Println("Returns:")
	fmt.Println(result)
}

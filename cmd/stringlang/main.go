package main

import (
	"fmt"
	"github.com/skius/stringlang"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("Usage: ./stringlang [<program.stringlang>] [..<args>]")
		return
	}
	if len(os.Args) == 1 {
		r := new(Repl)
		r.Run()
		return
	}

	file := os.Args[1]
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	expr, err := stringlang.Parse(content)
	if err != nil {
		panic(err)
	}

	result, err := evalOrTimeout(exampleContext(true), expr, time.Second*30)
	if err != nil {
		panic(err)
	}

	fmt.Println("Returns:")
	fmt.Println(strings.ReplaceAll(result, `\n`, "\n"))
}

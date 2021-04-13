package main

import (
	"bufio"
	"fmt"
	"github.com/skius/stringlang"
	"github.com/skius/stringlang/ast"
	"os"
	"strings"
	"time"
)

func repl() {
	fmt.Println("Welcome to the StringLang REPL!")
	fmt.Println("Enter as much code as you want, run it by ending your input with ';;', repeat!")
	fmt.Println("Reset your program using 'reset;;' and quit the REPL using 'quit;;'.")

	program := ast.EmptyProgram()

	for {
		input := readLines()

		// Special keywords
		if strings.HasSuffix(input, "reset;;") {
			fmt.Println("Resetting REPL... Reset!")
			program = ast.EmptyProgram()
			continue
		}

		if strings.HasSuffix(input, "quit;;") {
			break
		}

		// Remove trailing whitespace and semicolon, because StringLang programs end without semicolon
		source := strings.Trim(input, "\t\n\r ;")

		newExpr, err := stringlang.Parse([]byte(source))
		if err != nil {
			fmt.Println("There was an error parsing your program: ", err)
			continue
		}

		// Combine the program extension with the old program
		newProgram := newExpr.(ast.Program)
		program.Funcs = append(program.Funcs, newProgram.Funcs...)
		program.Code = append(program.Code, newProgram.Code...)


		result, err := evalOrTimeout(program, time.Second * 30)
		if err != nil {
			fmt.Println("There was an error running your program: ", err)
			continue
		}
		fmt.Println(result)
	}

	fmt.Println("Exiting REPL.")
}

func readLines() string {
	in    := bufio.NewReader(os.Stdin)
	input := ""

	// Read in lines of StringLang until ;;
	for !strings.HasSuffix(input, ";;") {
		temp, err := in.ReadString('\n')
		if err != nil {
			panic(err)
		}
		temp = strings.TrimSpace(temp)
		input += " " + temp
	}
	return strings.TrimSpace(input)
}

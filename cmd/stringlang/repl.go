package main

import (
	"bufio"
	"fmt"
	"github.com/skius/stringlang"
	"github.com/skius/stringlang/ast"
	"github.com/skius/stringlang/errors"
	"github.com/skius/stringlang/token"
	"os"
	"runtime"
	"strings"
	"time"
)


var Reset  = "\033[0m"
var Red    = "\033[31m"
var Green  = "\033[32m"
var Yellow = "\033[33m"
var Blue   = "\033[34m"
var Purple = "\033[35m"
var Cyan   = "\033[36m"
var Gray   = "\033[37m"
var White  = "\033[97m"

func repl() {
	if runtime.GOOS == "windows" {
		Reset  = ""
		Red    = ""
		Green  = ""
		Yellow = ""
		Blue   = ""
		Purple = ""
		Cyan   = ""
		Gray   = ""
		White  = ""
	}

	fmt.Println("Welcome to the StringLang REPL!")
	fmt.Println("Enter code, run it by pressing ENTER, repeat!")
	fmt.Println("Note: This REPL stores all code you've entered and will execute everything each time you enter new code.")
	fmt.Println("Reset your program using 'reset;;' and quit the REPL using 'quit;;'.")

	program := ast.EmptyProgram()
	stdin := bufio.NewReader(os.Stdin)

	repl:
	for {
		input := ""
		var newExpr stringlang.Expr
		for {
			temp, err := stdin.ReadString('\n')
			if err != nil {
				panic(err)
			}

			// Special keywords
			trimmedTemp := strings.TrimSpace(temp)
			if strings.HasSuffix(trimmedTemp, "reset;;") {
				fmt.Println("Resetting REPL... Reset!")
				program = ast.EmptyProgram()
				continue repl
			}
			if strings.HasSuffix(trimmedTemp, "quit;;") {
				break repl
			}

			input += temp

			newExpr, err = stringlang.Parse([]byte(input))
			if err == nil {
				// newExpr got successfully parsed into a program, let's execute it.
				break
			}
			pErr := err.(*errors.Error)
			if pErr.ErrorToken.Type == token.EOF {
				// If err is of the form "expected <something>; got: end-of-file", we know the program was potentially
				// correct, just incomplete, so we should keep reading.
				continue
			}
			unexpectedToken := pErr.ErrorToken.Lit
			if strings.HasPrefix(string(unexpectedToken), `"`) {
				// Start of a multiline string, keep reading.
				continue
			}

			// Otherwise there's no chance the program could become correct, so we have to reset
			fmt.Println("There was an error parsing your input: ", err)
			input = ""
		}

		// Combine the program extension with the old program
		newProgram := newExpr.(ast.Program)
		program.Funcs = append(program.Funcs, newProgram.Funcs...)
		program.Code = append(program.Code, newProgram.Code...)

		if len(newProgram.Code) == 0 {
			// No new top-level code, so no need to run anything
			continue
		}

		result, err := evalOrTimeout(program, time.Second * 30)
		if err != nil {
			fmt.Println("There was an error running your program: ", err)
			continue
		}
		fmt.Println(Yellow + result + Reset)
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

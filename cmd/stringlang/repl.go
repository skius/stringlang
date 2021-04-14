package main

import (
	"bufio"
	"fmt"
	"github.com/skius/stringlang"
	"github.com/skius/stringlang/ast"
	"github.com/skius/stringlang/errors"
	"github.com/skius/stringlang/token"
	"golang.org/x/term"
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

	// _ is missing setIndent
	printLn, readLn, cleanup, _ := initTerminal()
	defer cleanup()

	printLn("Welcome to the StringLang REPL!")
	printLn("Enter code, run it by pressing ENTER, repeat!")
	printLn("Note: This REPL stores all code you've entered and will execute everything each time you enter new code.")
	printLn("Reset your program using 'reset;;' and quit the REPL using 'quit;;'.")

	program := ast.EmptyProgram()

	repl:
	for {
		indentLevel := 0
		input := ""
		var newExpr stringlang.Expr
		// Read-parse loop until we either know input can't be valid a StringLang fragment, or we have a successful parse
		for {
			// Uncomment for automated basic and bad indentation support
			// setIndent(indentLevel)
			temp := readLn()

			// Special keywords
			trimmedTemp := strings.TrimSpace(temp)
			if strings.HasSuffix(trimmedTemp, "reset;;") {
				printLn("Resetting REPL... Reset!")
				program = ast.EmptyProgram()
				continue repl
			}
			if strings.HasSuffix(trimmedTemp, "quit;;") {
				break repl
			}

			// Check if we need to change indentation
			opens := strings.Count(temp, "{")
			closes := strings.Count(temp, "}")
			indentLevel = indentLevel + opens - closes

			input += temp

			var err error
			newExpr, err = stringlang.Parse([]byte(input))
			if err == nil {
				// newExpr got successfully parsed into a program, let's execute it
				break
			}
			pErr := err.(*errors.Error)
			if pErr.ErrorToken.Type == token.EOF {
				// If err is of the form "expected <something>; got: end-of-file", we know the program was potentially
				// correct, just incomplete, so we should keep reading
				continue
			}
			unexpectedToken := pErr.ErrorToken.Lit
			if strings.HasPrefix(string(unexpectedToken), `"`) {
				// Start of a multiline string, keep reading
				continue
			}

			// Otherwise there's no chance the program could become correct, so we have to reset
			printLn("There was an error parsing your input: ", err)
			input = ""
			indentLevel = 0
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
			printLn("There was an error running your program: ", err)
			continue
		}
		printLn(Yellow + result + Reset)

	}

	printLn("Exiting REPL.")
}

func initTerminal() (func(...interface{}), func() string, func(), func(int)) {
	var printLn func(...interface{})
	var readLn func() string
	var cleanup = func() {}
	var setIndent func(int)


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

		var indent int
		setIndent = func(i int) {
			indent = i
		}

		stdin := bufio.NewReader(os.Stdin)
		printLn = func(a ...interface{}) {
			fmt.Println(a...)
		}
		readLn = func() string {
			fmt.Print(genSpaces(indent * 4))
			s, err := stdin.ReadString('\n')
			if err != nil {
				panic(err)
			}
			// ReadString returns line incl. delimiter = line-break, so keep 's' as is
			return s
		}
	} else {
		// Can use proper terminal
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}

		t := term.NewTerminal(os.Stdin, "> ")

		setIndent = func(i int) {
			// Need to add "> " if we want it even with indentation
			t.SetPrompt(genSpaces(i * 4))
		}

		printLn = func(a ...interface{}) {
			_, err := t.Write([]byte(fmt.Sprintln(a...)))
			if err != nil {
				panic(err)
			}
		}
		readLn = func() string {
			s, err := t.ReadLine()
			if err != nil {
				panic(err)
			}
			// Terminal.ReadLine returns line without line-break, so let's add it to s
			return s + "\n"
		}
		cleanup = func() {
			term.Restore(int(os.Stdin.Fd()), oldState)
		}
	}
	return printLn, readLn, cleanup, setIndent
}

func genSpaces(i int) string {
	s := ""
	for j := 0; j < i; j++ {
		s += " "
	}
	return s
}

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

	// TODO: Figure out better what to do when stack overflows (because subsequent calls will just return "" because stack will still be overflown
	// TODO: Figure out why in an infinite loop Ctrl-C doesnt work? only after 30sec timeout, then EOF

	// _ is missing setIndent
	printLn, readLn, cleanup, _, setMultiLine := initTerminal()
	defer cleanup()

	printLn("Welcome to the StringLang REPL!")
	printLn("Enter code, run it by pressing ENTER, repeat!")
	printLn("The special variable '_' can be used to refer to the previous result.")
	printLn("Reset your program using 'reset;;' and quit the REPL using 'quit;;'.")

	funcs := []ast.FuncDecl{}
	context := exampleContext()

	repl:
	for {
		indentLevel := 0
		input := ""
		var newExpr stringlang.Expr
		// Read-parse loop until we either know input can't be valid a StringLang fragment, or we have a successful parse
		for {
			if input == "" {
				// input can only ever be "" if this is the first iteration of
				// the read-parse loop ~=~ the beginning of a new expression
				setMultiLine(false)
			} else {
				// We only reach the second iteration if the first ended unexpectedly but is recoverable, i.e. multiline
				setMultiLine(true)
			}

			// Uncomment for automated basic and bad indentation support
			// setIndent(indentLevel)

			temp := readLn()

			// Special keywords
			trimmedTemp := strings.TrimSpace(temp)
			if strings.HasSuffix(trimmedTemp, "reset;;") {
				printLn("Resetting REPL... Reset!")
				funcs = []ast.FuncDecl{}
				context = exampleContext()
				continue repl
			}
			if strings.HasSuffix(trimmedTemp, "quit;;") {
				break repl
			}

			// Basic and easily broken check if we need to change indentation
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

			// Otherwise there's no chance the program could become correct, so we have to reset this expression
			printLn("There was an error parsing your input: ", err)
			input = ""
			indentLevel = 0
		}

		newProgram := newExpr.(ast.Program)

		// Combine the new functions with the old ones
		// Order is important, newer functions need to be at the end of the slice
		// such that they properly replace earlier definitions
		funcs = append(funcs, newProgram.Funcs...)

		newProgram.Funcs = funcs

		if len(newProgram.Code) == 0 {
			// No new top-level code, so no need to run anything
			continue
		}

		// Eval by reusing context, so we store previous computations
		result, err := evalOrTimeout(context, newProgram, time.Second * 30)
		if err != nil {
			printLn("There was an error running your program: ", err)
			continue
		}
		// Update special variable '_' to refer to result
		context.VariableMap["_"] = ast.Val(result)
		printLn(Yellow + result + Reset)

	}

	printLn("Exiting REPL.")
}

func initTerminal() (func(...interface{}), func() string, func(), func(int), func(bool)) {
	var printLn func(...interface{})
	var readLn func() string
	var cleanup = func() {}
	var setIndent func(int)
	var setMultiLine func(bool)

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

		// currently unused
		var multiLine bool
		setMultiLine = func(is bool) {
			multiLine = is
		}

		var indent int
		setIndent = func(i int) {
			indent = i
		}

		stdin := bufio.NewReader(os.Stdin)
		printLn = func(a ...interface{}) {
			fmt.Println(a...)
		}
		readLn = func() string {
			fmt.Print("> " + genSpaces(indent * 4))
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

		var _updatePrompt func()

		var multiLine bool
		setMultiLine = func(is bool) {
			multiLine = is
			_updatePrompt()
		}

		var indent int
		setIndent = func(i int) {
			indent = i
			_updatePrompt()
		}

		_getMultiLinePrefix := func() string {
			if multiLine {
				return "  "
			} else {
				return "> "
			}
		}

		_getIndentPrefix := func() string {
			return genSpaces(indent * 4)
		}

		_updatePrompt = func() {
			t.SetPrompt(_getMultiLinePrefix() + _getIndentPrefix())
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
	return printLn, readLn, cleanup, setIndent, setMultiLine
}

func genSpaces(i int) string {
	s := ""
	for j := 0; j < i; j++ {
		s += " "
	}
	return s
}

package repl

import (
	"github.com/skius/stringlang"
	"github.com/skius/stringlang/ast"
	"github.com/skius/stringlang/internal/frontend/errors"
	"github.com/skius/stringlang/internal/frontend/token"
	"strings"
	"time"
)

type Repl struct {
	Context      *stringlang.Context
	UserFuncs    []ast.FuncDecl
	IndentLevel  int
	PartialParse string
	T            Terminal
}

func Init(t Terminal) *Repl {
	r := new(Repl)
	r.Context = stringlang.ExampleContext(false)
	r.UserFuncs = make([]ast.FuncDecl, 0)
	r.IndentLevel = 0
	r.PartialParse = ""
	r.T = t
	return r
}

func (r *Repl) PrintWelcome() {
	r.T.PrintLn("Welcome to the StringLang REPL!")
	r.T.PrintLn("Enter code, run it by pressing ENTER, repeat!")
	r.T.PrintLn("The special variable '_' can be used to refer to the previous result.")
	r.T.PrintLn("Reset your program using 'reset;;' and quit the REPL using 'quit;;' or Ctrl-C.")
}

func (r *Repl) ResetPartial() {
	r.IndentLevel = 0
	r.PartialParse = ""
}

func (r *Repl) FullReset() {
	r.T.PrintLn("Resetting REPL... Reset!")
	r.UserFuncs = []ast.FuncDecl{}
	r.Context = stringlang.ExampleContext(false)
	r.ResetPartial()
}

func (r *Repl) IsNewPartialParse() bool {
	// The partial parse can only ever be "" if this is the first iteration of
	// the read-parse loop ~=~ the beginning of a new expression
	return r.PartialParse == ""
}

func (r *Repl) UpdateIndent(s string) {
	// Basic check if we need to change indentation
	opens := strings.Count(s, "{")
	closes := strings.Count(s, "}")
	r.IndentLevel = r.IndentLevel + opens - closes
	r.T.SetIndent(r.IndentLevel)
}

func (r *Repl) SetSpecial(val string) {
	r.Context.VariableMap["_"] = ast.Val(val)
}

func (r *Repl) ReadExpr() (expr stringlang.Expr, reset bool, quit bool) {
	// Clean up previous work
	r.ResetPartial()
	// Read-parse loop until we either know input can't be valid a StringLang fragment, or we have a successful parse
	for {
		if r.IsNewPartialParse() {
			r.T.SetMultiLine(false)
		} else {
			// We only reach the second iteration if the first ended unexpectedly but is recoverable,
			// i.e. non-empty PartialParse -> multiline
			r.T.SetMultiLine(true)
		}

		line := r.T.ReadLn()

		// Special keywords
		trimmedTemp := strings.TrimSpace(line)
		if isCmd(trimmedTemp, "reset") {
			return nil, true, false
		}
		if isCmd(trimmedTemp, "quit") {
			return nil, false, true
		}

		r.UpdateIndent(line)

		r.PartialParse += line

		var err error
		expr, err = stringlang.Parse([]byte(r.PartialParse))
		if err == nil {
			// newExpr got successfully parsed into a program, let's execute it
			return expr, false, false
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
		r.T.PrintLn("There was an error parsing your input: ", err)
		r.ResetPartial()
	}
}

func (r *Repl) Run() {
	t := r.T
	defer t.Cleanup()
	r.PrintWelcome()

	for {
		expr, reset, quit := r.ReadExpr()
		if quit {
			break
		}
		if reset {
			r.FullReset()
			continue
		}

		prog := expr.(ast.Program)
		if len(prog.Code) == 0 {
			// No new top-level code, so no need to run anything
			// Store new user functions, however
			for i := range prog.Funcs {
				fn := prog.Funcs[i]
				r.Context.UserFunctionMap[fn.Identifier] = fn
			}
			continue
		}

		// Eval by reusing context, so we store previous computations
		result, err := stringlang.EvalOrTimeout(r.Context, prog, time.Second*5)
		if err != nil {
			t.PrintLn("There was an error running your program: ", err)
			continue
		}
		// Update special variable '_' to refer to result
		r.SetSpecial(result)
		t.PrintLn(t.Color(Yellow) + result + t.ResetColor())
	}

	t.PrintLn("Exiting REPL.")
}

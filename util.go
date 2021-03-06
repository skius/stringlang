package stringlang

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func EvalOrTimeout(ctx *Context, expr Expr, timeout time.Duration) (string, error) {
	exit := ctx.GetExitChannel()

	resultChan := make(chan string)
	errChan := make(chan interface{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- r
			}
		}()
		resultChan <- string(expr.Eval(ctx))
	}()

	var result string
	select {
	case result = <-resultChan:
	case err := <-errChan:
		return "", errors.New(fmt.Sprint(err))
	case <-time.After(timeout):
		exit <- 1
		return "", errors.New(fmt.Sprint("Program timed out after ", timeout))
	}
	return result, nil
}

func ExampleContext(limitStack bool) *Context {
	rand.Seed(time.Now().UnixNano())
	funcs := map[string]func([]string) string{
		"random": func(args []string) string {
			num := len(args)
			if num == 0 {
				return strconv.Itoa(rand.Intn(10) + 1)
			} else if num == 1 {
				val, err := strconv.Atoi(args[0])
				if err == nil {
					return strconv.Itoa(rand.Intn(val) + 1)
				}
			}
			return args[rand.Intn(num)]
		},
		"length": func(args []string) string {
			return strconv.Itoa(len(args[0]))
		},
	}

	args := make([]string, len(flag.Args()))
	copy(args, flag.Args())

	ctx := NewContext(args, funcs)
	// Add special eval function, needs reference to ctx to work
	// Modifies context, i.e. eval adds support for unhygienic macros
	// Only works if ctx gets reused, like what REPL is doing
	// TODO: eval in a function body actually evals using the parent scope. no good.
	ctx.FunctionMap["eval"] = func(args []string) string {
		if len(args) == 0 {
			return ""
		}
		src := args[0]
		expr, err := Parse([]byte(src))
		if err != nil {
			return ""
		}
		result := expr.Eval(ctx)
		return string(result)
	}
	if limitStack {
		ctx.SetMaxStackSize(100 * 1024 * 1024) // 100MB limit for programs
	}

	return ctx
}

package main

import (
	"errors"
	"fmt"
	"github.com/skius/stringlang"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func genSpaces(i int) string {
	s := ""
	for j := 0; j < i; j++ {
		s += " "
	}
	return s
}

func isCmd(s, cmd string) bool {
	return strings.HasSuffix(s, cmd+";;")
}

func evalOrTimeout(ctx *stringlang.Context, expr stringlang.Expr, timeout time.Duration) (string, error) {
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
	case err := <- errChan:
		return "", errors.New(fmt.Sprint(err))
	case <-time.After(timeout):
		exit <- 1
		return "", errors.New(fmt.Sprint("Program timed out after ", timeout))
	}
	return result, nil
}

func exampleContext(limitStack bool) *stringlang.Context {
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

	// We want the .stringlang file, and not the interpreter's path to be argument 0
	args := make([]string, len(os.Args)-1)
	for i := 1; i < len(os.Args); i++ {
		args[i-1] = os.Args[i]
	}

	ctx := stringlang.NewContext(args, funcs)
	if limitStack {
		ctx.SetMaxStackSize(100 * 1024 * 1024) // 100MB limit for programs
	}
	return ctx
}

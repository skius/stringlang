package main

import (
	"errors"
	"fmt"
	"github.com/skius/stringlang"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("Usage: ./stringlang [<program.stringlang>] [..<args>]")
		return
	}
	if len(os.Args) == 1 {
		repl()
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

	result, err := evalOrTimeout(expr, time.Second * 30)
	if err != nil {
		panic(err)
	}

	fmt.Println("Returns:")
	fmt.Println(strings.ReplaceAll(result, `\n`, "\n"))
}

func evalOrTimeout(expr stringlang.Expr, timeout time.Duration) (string, error) {
	ctx := exampleContext()
	exit := ctx.GetExitChannel()

	resultChan := make(chan string)
	go func() {
		resultChan <- string(expr.Eval(ctx))
	}()

	var result string
	select {
	case result = <- resultChan:
	case <-time.After(timeout):
		exit <- 1
		return "", errors.New(fmt.Sprint("Program timed out after ", timeout))
	}
	return result, nil
}

func exampleContext() *stringlang.Context {
	rand.Seed(time.Now().UnixNano())
	funcs := map[string]func([]string)string {
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
		"length" : func(args []string) string {
			return strconv.Itoa(len(args[0]))
		},
	}

	// We want the .stringlang file, and not the interpreter's path to be argument 0
	args := make([]string, len(os.Args) - 1)
	for i := 1; i < len(os.Args); i++ {
		args[i-1] = os.Args[i]
	}

	ctx := stringlang.NewContext(args, funcs)
	ctx.SetMaxStackSize(100 * 1024 * 1024) // 100MB limit for programs
	return ctx
}

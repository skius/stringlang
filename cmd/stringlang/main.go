package main

import (
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
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./stringlang <program.stringlang> ..<args>")
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
	ctx.MaxWhileIter = 10000
	ctx.MaxStackSize = 100 * 1024 * 1024 // 100MB limit for variables

	result := string(expr.Eval(ctx))

	fmt.Println("Returns:")
	fmt.Println(strings.ReplaceAll(result, `\n`, "\n"))
}
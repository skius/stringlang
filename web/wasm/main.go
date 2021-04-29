package main

import (
	"fmt"
	"github.com/skius/stringlang"
	"syscall/js"
)

var ctx *stringlang.Context

func Eval(s string) string {
	// TODO: Extract the repl from cmd/stringlang into a separate package and use it here
	expr, err := stringlang.Parse([]byte(s))
	if err != nil {
		fmt.Println(err)
	}

	res := expr.Eval(ctx)

	return string(res)
}

func wrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		input := args[0].String()
		res := Eval(input)
		return res
	})
	return jsonFunc
}

func main() {
	fmt.Println("Running!")

	ctx = stringlang.NewContext([]string{}, map[string]func([]string) string{})

	js.Global().Set("stringlang", wrapper())
	<-make(chan bool)
}

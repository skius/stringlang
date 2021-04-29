package main

import (
	"fmt"
	"github.com/skius/stringlang/cmd/stringlang/repl"
	"syscall/js"
)

type wasm struct {
	inChan chan<- string
}

// Receives input from wasm and forwards it to the JSTerminal
func (w *wasm) input(s string) {
	w.inChan <- s
}

func wrapperInput(w *wasm) js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		input := args[0].String()
		w.input(input)
		return nil
	})
	return jsonFunc
}

func main() {
	fmt.Println("Running!")

	w := new(wasm)
	inCh := make(chan string, 1)
	w.inChan = inCh

	js.Global().Set("SLInput", wrapperInput(w))

	t := NewJSTerminal(inCh)
	r := repl.Init(t)

	r.Run()

	<-make(chan bool)
}

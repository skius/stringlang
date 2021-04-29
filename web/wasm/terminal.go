package main

import (
	"fmt"
	"syscall/js"
)

type JSTerminal struct {
	print    func(string)
	readChan <-chan string
}

func NewJSTerminal(in chan string) *JSTerminal {
	t := new(JSTerminal)
	t.readChan = in
	t.print = func(s string) {
		js.Global().Call("write", s)
	}
	return t
}

func (t *JSTerminal) PrintLn(i ...interface{}) {
	t.print(fmt.Sprintln(i...))
}

func (t *JSTerminal) ReadLn() string {
	return <-t.readChan
}

func (t *JSTerminal) SetIndent(i int) {

}

func (t *JSTerminal) SetMultiLine(b bool) {

}

func (t *JSTerminal) PrintPrompt() {

}

func (t *JSTerminal) Cleanup() {

}

func (t *JSTerminal) Color(i int) string {
	return ""
}

func (t *JSTerminal) ResetColor() string {
	return ""
}

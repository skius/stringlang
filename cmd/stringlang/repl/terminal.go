package repl

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"runtime"
)

const (
	Red int = iota
	Green
	Yellow
	Blue
	Purple
	Cyan
	Gray
	White
)

const SpacesPerIndent = 4

type Terminal interface {
	PrintLn(...interface{})
	ReadLn() string

	SetIndent(int)
	SetMultiLine(bool)
	PrintPrompt()
	Cleanup()

	Color(int) string
	ResetColor() string
}

func DefaultTerminal() Terminal {
	if runtime.GOOS == "windows" {
		t := new(SimpleTerminal)
		t.Init(bufio.NewReader(os.Stdin))
		return t
	} else {
		t := new(UnixTerminal)
		t.Init(os.Stdin)
		return t
	}
}

/*
	Terminal for Windows etc.
*/

type SimpleTerminal struct {
	multiLine   bool
	indentLevel int
	in          *bufio.Reader
}

func (t *SimpleTerminal) Init(reader *bufio.Reader) {
	t.in = reader
}

func (t *SimpleTerminal) PrintLn(a ...interface{}) {
	fmt.Println(a...)
}
func (t *SimpleTerminal) ReadLn() string {
	t.PrintPrompt()
	s, err := t.in.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return s
}

func (t *SimpleTerminal) SetIndent(i int) {
	t.indentLevel = i
}
func (t *SimpleTerminal) getIndentPrefix() string {
	return genSpaces(SpacesPerIndent * t.indentLevel)
}
func (t *SimpleTerminal) SetMultiLine(multiLine bool) {
	t.multiLine = multiLine
}
func (t *SimpleTerminal) getMultiLinePrefix() string {
	if t.multiLine {
		return "  "
	} else {
		return "> "
	}
}
func (t *SimpleTerminal) PrintPrompt() {
	fmt.Print(t.getMultiLinePrefix() + t.getIndentPrefix())
}

func (*SimpleTerminal) Cleanup() {}
func (*SimpleTerminal) Color(int) string {
	return ""
}
func (*SimpleTerminal) ResetColor() string {
	return ""
}

/*
	Terminal for Unix
*/

type UnixTerminal struct {
	multiLine   bool
	indentLevel int
	t           *term.Terminal
	oldState    *term.State
	inout       *os.File
	readChan    chan string
}

func (t *UnixTerminal) Init(inout *os.File) {
	oldState, err := term.MakeRaw(int(inout.Fd()))
	if err != nil {
		panic(err)
	}
	t.oldState = oldState
	t.t = term.NewTerminal(inout, "")
	t.inout = inout

	t.readChan = make(chan string)
	go func() {
		defer t.Cleanup()
		for {
			s, err := t.t.ReadLine()
			if err == io.EOF {
				t.Cleanup()
				fmt.Println("\nExiting.")
				os.Exit(1)
			}
			if err != nil {
				panic(err)
			}
			// Terminal.ReadLine returns line without line-break, so let's add it to s
			s += "\n"
			t.readChan <- s
		}
	}()
}

func (t *UnixTerminal) PrintLn(a ...interface{}) {
	_, err := t.t.Write([]byte(fmt.Sprintln(a...)))
	if err != nil {
		panic(err)
	}
}
func (t *UnixTerminal) ReadLn() string {
	t.PrintPrompt()
	s := <-t.readChan
	return s
}

func (t *UnixTerminal) SetIndent(i int) {
	t.indentLevel = i
}
func (t *UnixTerminal) getIndentPrefix() string {
	return genSpaces(SpacesPerIndent * t.indentLevel)
}
func (t *UnixTerminal) SetMultiLine(multiLine bool) {
	t.multiLine = multiLine
}
func (t *UnixTerminal) getMultiLinePrefix() string {
	if t.multiLine {
		return "  "
	} else {
		return "> "
	}
}
func (t *UnixTerminal) PrintPrompt() {
	_, err := t.t.Write([]byte(t.getMultiLinePrefix() + t.getIndentPrefix()))
	if err != nil {
		panic(err)
	}
}

func (t *UnixTerminal) Cleanup() {
	err := term.Restore(int(t.inout.Fd()), t.oldState)
	if err != nil {
		fmt.Println(err)
	}
}
func (*UnixTerminal) Color(c int) string {
	switch c {
	case Red:
		return "\033[31m"
	case Green:
		return "\033[32m"
	case Yellow:
		return "\033[33m"
	case Blue:
		return "\033[34m"
	case Purple:
		return "\033[35m"
	case Cyan:
		return "\033[36m"
	case Gray:
		return "\033[37m"
	case White:
		return "\033[97m"
	default:
		return ""
	}
}
func (*UnixTerminal) ResetColor() string {
	return "\033[0m"
}

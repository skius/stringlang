package main

import (
	"bufio"
	"fmt"
	"github.com/skius/stringlang"
	"os"
	"strings"
	"time"
)

func repl() {
	fmt.Println("Welcome to the StringLang* REPL! (* function definitions not supported)")
	fmt.Println("Enter as much code as you want, run it by ending your input with ';;', repeat!")
	fmt.Println("Reset your program using 'reset;;' and quit the REPL using 'quit;;'.")

	program := ""

	for {
		input := readLines()

		// Special keywords
		if input == "reset;;" {
			fmt.Println("Resetting REPL... Reset!")
			program = ""
			continue
		}

		if input == "quit;;" {
			break
		}

		// Remove trailing whitespace and semicolon, because StringLang programs end without semicolon
		code := strings.Trim(input, "\t\n\r ;")
		// We keep track of the full program as its source code; a better way would be extending the AST
		if program == "" {
			program = code
		} else {
			program += "; " + code
		}

		expr, err := stringlang.Parse([]byte(program))
		if err != nil {
			fmt.Println("There was an error parsing your program: ", err)
			fmt.Println("Resetting REPL... Reset!")
			program = ""
			continue
		}

		result, err := evalOrTimeout(expr, time.Second * 30)
		if err != nil {
			fmt.Println("There was an error running your program: ", err)
			fmt.Println("Resetting REPL... Reset!")
			program = ""
			continue
		}
		fmt.Println(result)
	}

	fmt.Println("Exiting REPL.")
}

func readLines() string {
	in    := bufio.NewReader(os.Stdin)
	input := ""

	// Read in lines of StringLang until ;;
	for !strings.HasSuffix(input, ";;") {
		temp, err := in.ReadString('\n')
		if err != nil {
			panic(err)
		}
		temp = strings.TrimSpace(temp)
		input += " " + temp
	}
	return strings.TrimSpace(input)
}

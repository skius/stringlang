package repl

import "strings"

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

package ast

import "github.com/skius/stringlang/token"

func attribToString(a Attrib) string {
	return string(a.(*token.Token).Lit)
}
func unescape(s string) string {
	in := []rune(s)
	out := make([]rune, 0, len(in))
	var escape bool
	for _, r := range in {
		switch {
		case escape:
			switch r {
			case 'n':
				out = append(out, '\n')
			default:
				out = append(out, r)
			}
			escape = false
		case r == '\\':
			escape = true
		default:
			out = append(out, r)
		}
	}
	return string(out)
}
func boolOf(v Val) bool {
	return v != "false" && v != ""
}
func CheckSize(m map[Var]Val) (total int64) {
	for k, v := range m {
		total += int64(len(k)) + int64(len(v))
	}
	return
}

const (
	SigExternalExit = iota + 1
	SigOutOfMemory
)

// checkExit returns true if we need to exit
func checkExit(c *Context) bool {
	if c.limitStackSize && CheckSize(c.VariableMap) > c.MaxStackSize {
		select {
		case c.exitChannel <- SigOutOfMemory:
		default:
		}
		return true
	}
	select {
	case <-c.exitChannel:
		// c.exitChannel <- sig // No need to propagate I think?
		return true
	default:
		return false
	}
}

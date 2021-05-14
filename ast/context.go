package ast

func NewContext(args []string, funcs map[string]func([]string) string, parseFn func([]byte) (Expr, error)) *Context {
	return &Context{
		Args:            args,
		VariableMap:     make(map[Var]Val),
		UserFunctionMap: make(map[string]FuncDecl),
		FunctionMap:     funcs,
		MaxStackSize:    -1,
		limitStackSize:  false,
		exitChannel:     make(chan int, 1),
		parseFn:         parseFn,
	}
}

type Context struct {
	Args            []string
	VariableMap     map[Var]Val
	FunctionMap     map[string]func([]string) string
	UserFunctionMap map[string]FuncDecl
	MaxStackSize    int64
	exitChannel     chan int
	limitStackSize  bool
	parseFn         func([]byte) (Expr, error) /* Ugly hack to avoid illegal circular imports */
}

func (c *Context) SetMaxStackSize(sz int64) {
	c.MaxStackSize = sz
	c.limitStackSize = true
	if sz < 0 {
		c.limitStackSize = false
	}
}

func (c *Context) GetExitChannel() chan int {
	if c.exitChannel == nil {
		c.exitChannel = make(chan int, 1)
	}
	return c.exitChannel
}

func (c *Context) FuncNames() Set {
	names := make(Set, len(c.FunctionMap) + len(c.UserFunctionMap))
	for fId := range c.UserFunctionMap {
		names.Add(fId)
	}
	for fId := range c.FunctionMap {
		names.Add(fId)
	}
	return names
}

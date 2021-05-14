# StringLang

An interpreted, expression-oriented language where everything evaluates to strings.

## Usage

### Installation

You need Go to build the official interpreter from this repository.  
Run `go get github.com/skius/stringlang/cmd/stringlang` to install the CLI on your machine.

### CLI

Run StringLang programs using `stringlang <program.stringlang> <arg0> <arg1> ...`,
or alternatively run the StringLang REPL by running `stringlang` with no arguments.

### Running from code

To interpret StringLang code from your Go program, all you need is the following:

```go
package main

import (
   "github.com/skius/stringlang"
   "os"
   "strconv"
)

func main() {
   source := `"Replace me with the " + "source code of your StringLang program"`
   expr, err := stringlang.Parse([]byte(source))
   if err != nil {
      panic(err)
   }
   // The built-in functions your StringLang program will have access to
   funcs := map[string]func([]string) string{
      "length": func(as []string) string { return strconv.Itoa(len(as[0])) },
      // Add more built-in functions here
   }
   // Arguments to your StringLang program
   args := os.Args
   ctx := stringlang.NewContext(args, funcs)
   result := expr.Eval(ctx)
}
```

See the CLI's [main.go](cmd/stringlang/main.go) for a more advanced example.

### Contributing

Feel free to open Issues and Pull Requests! The language specification and interpreter is by no means final.
To change the syntax of the language, you'll need to modify `lang.bnf` and also `ast/ast.go` if you need new structures.
Every time you change `lang.bnf`, you need to run `gocc -o "./internal/frontend" -p "github.com/skius/internal/frontend" lang.bnf` to generate the new parser and lexer.

Get the parser generator `gocc` used in this project from: [goccmack/gocc](https://github.com/goccmack/gocc)

### Example StringLang programs

See [`stringlang_programs/`](stringlang_programs)

## Syntax

The following code block contains a simple (and not 100% correct, see Caveats section below the block) overview of what
constitutes a `StringLang` source file. The complete grammar can be found in `lang.bnf`.

```
identifier: any string of alphanumeric and underscore characters not beginning with a digit
number: non-negative integer
string_literal: a "double-quoted" string, containing any sequence of characters, escapes allowed using '\'
comment: any text enclosed by /* and */, except '/*' and '*/' (no nested comments)


program:
     header block           

header:
    function1 function2 ... functionN                       // N may be 0  

function:
    fun identifier(param1, param2, ..., paramN) { block }   // paramX are identifiers, N may be 0

block:
    expression1; expression2; ...; expressionN              // no trailing semi-colon and N cannot be 0

expression:
    string_literal
    identifier
    identifier = expression
    $number
    expression[expression]
    expression[number]
    expression(expression1, expression2, ..., expressionN)  // N can be 0
    
    expression || expression
    expression && expression
    expression != expression
    expression == expression
    expression + expression
    
    (expression)
    
    ifelse
    while
    
    lambda
    
ifelse:
    if (expression) { block } else { block }
    if (expression) { block } else ifelse                   // This effectively allows else-if 
    
while:
    while (expression) { block }

lambda:
    fun(param1, param2, ..., paramN) { block }              // paramX are identifiers, N may be 0
```

### Caveats
While the above description gives a good overview of `StringLang`, there are some important notes to be made:
1) `identifier = expression` may only appear (but possibly repeated) at the beginning of an expression without parentheses, e.g. 
   `a = b = c = "foo";`. Otherwise they need to be inside parentheses, e.g. `"foo" + (a = c = "bar")`.
2) A block is a non-empty, `;` -separated list of expressions, with no trailing `;`.
3) A call (`identifier(expr1, ..., exprN)`) may also be without arguments, i.e. `identifier()`

## Semantics

### Values

All expressions/values in `StringLang` are Strings.

### Functions

Functions in `StringLang` can either be user-defined (i.e. they are written in `StringLang` and reside
as a top-level function in the header of the interpreted file, or are lambdas), or they can be built-in 
(i.e. supplied to the interpreter via the [context](#context-functions-and-arguments)). 
Because built-in functions are written in Go, they can accomplish anything Go can accomplish. 

User-defined `StringLang` functions are evaluated with their own completely separate variable scope and
are mutually recursive. The arguments passed to them are bound to the variables in the corresponding parameter lists.

All functions in `StringLang` are pass-by-value, hence also strict.

#### Calls

In a `expr(args...)` call, finding the function corresponding to `expr` has the following order:
1. Check if `expr` is an `identifier` equal to some `f` and that there exists a *top-level, user-defined* function 
   `fun f(params) { block }`, if so, call that function.
2. Check if `expr` is an `identifier` equal to some `f` and that there exists a *built-in* 
   function that is mapped to that `f` in the provided [context](#context-functions-and-arguments), if so, call that.
3. If `expr` is not an `identifier` with corresponding named function, then [evaluate](#evaluation) `expr` and treat the
   resulting String as source code of a `lambda`, then call that lambda.

#### Lambdas

Lambdas in StringLang capture their environment by-value, and encode it as simple assignments at the top of the `block`.
Evaluating them (which is different to calling them) results in a canonical source representation of the lambda.  

Example: The expression `fun(x) { x + y }` is evaluated differently depending on its context. In particular,
the run-time value of the `identifier` `y` will be copied into the block, e.g. assuming `y` evaluates to `"value"`
in the current context the lambda becomes: `fun(x) { y = "value"; x + y }`, a context-insensitive lambda. Now it can
be evaluated into a String and returned: 
```
"fun(x) {
   y = \"value\";
   x + y
}"
```
At the site of a [call](#calls), this String can now be interpreted as the source code of a lambda, and called as such.  
Note that because we stored the value of `y` inside the lambda, the context in which we call this "lambda-turned-String"
does not matter.

Also note that lambdas are still user-defined functions, which are evaluated using a separate variable scope. Hence the
`y = "value"` in the block will not cause the caller's `y` to be set to `"value"`.

If you need more examples, see [lambdas.stringlang](stringlang_programs/lambdas.stringlang) or
[ski_combinator.stringlang](stringlang_programs/ski_combinator.stringlang) to see how the SK-calculus looks like in
StringLang.



### Binary Operators

The following list gives the binary operators, ordered from low to high precedence:
```
||   Or              Interprets operands as booleans           
&&   And             Interprets operands as booleans  
!=   Not Equals      Compares the values of the operands
==   Equals          Compares the values of the operands
+    Concatenation   Concatenates the operands
```
#### Note: 
1) The boolean value of a `StringLang` value is `true` if and only if the String is not `""` and 
   not `"false"`, else it is `false`.
2) The `StringLang` value of a boolean result of a logic operation is always either `"true"` or `"false"` 

### Evaluation 
```
string_literal ------------ The value of 'string_literal'
identifier ---------------- The current value of the variable with identifier 'identifier'.
identifier = expression --- The value of 'expression'. 
                            Side effect: Variable 'identifier' now has that value.
$number ------------------- The value of the 'number'-th (zero-indexed) argument to the program.
expr1[expr2 or number] ---- The character at position 'expr2' resp. 'number' of the value 
                            that 'expr1' evaluates to.
expr(expr1, ...) ---------- Function call to function 'expr' (See 'Functions' section) with arguments 'expr1, ...'.
                            Arguments are evaluated before passed to the function.
                            Evaluates to the function's return value.
expr1 <binop> expr2 ------- The value of the corresponding binary operation. 

ifelse -------------------- The value of the then-block if the condition evaluates to a true value,
                            or the value of the else-block if the condition evaluates to a false value.
while --------------------- The value of the last iteration of the block if it gets executed once, else "".
                            Side effects: All iterations might cause side effects
                            
lambda -------------------- The canonical source representation of the whole lambda as string, but with the blocks's
                            used variables captured by-value in the lambda.
                            
block --------------------- The value of the last expression in the block.
                            Side effects: Evaluates all expressions in the list.

```

#### Exceptions, errors, defaults

In the case of an invalid (out-of-bounds or NaN) character-access or argument expression,
the returned value is always `""`. The values of all variables are initialized to `""`.

### Context (functions and arguments)

To evaluate a `StringLang` program, the interpreter needs a `stringlang.Context` object.
It contains fields which allow the interpreter's user to supply their custom built-in functions and arguments to the program.
See `cmd/stringlang/main.go` for an example.

There is a channel available with `context.GetExitChannel()` for quickly killing the whole evaluation.
Additionally, one can set the (approximate) maximum stack space in bytes the `StringLang` program is allowed to use with
`context.SetMaxStackSize(int)` to a non-negative number. `cmd/stringlang/main.go` also contains examples
for these two features. 

## Why?

It started out as a simple String builder using provided arguments (read: `"my string".replace("$0", args[0])...`) 
for a chat bot, which allowed generating new commands in-chat. 
After switching to a parsed grammar for some simple operators, I figured why not just make it Turing-complete.
What could possibly go wrong allowing users to build their own commands...

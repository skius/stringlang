# StringLang

An interpreted, expression-oriented language where everything evaluates to strings.

## Examples

See `stringlang_programs/`

## Syntax

The following code block contains a simple (and not 100% correct, see Notes section below the block) overview of what
constitutes a `StringLang` source file. The complete grammar can be found in `lang.bnf`.

```
identifier: any string of alphanumeric and underscore characters not beginning with a digit
number: positive integer
string_literal: a "double-quoted" string, containing any sequence of characters


program = block

block:
    expression1; expression2; ...; expressionN              // Note: no trailing semi-colon and N cannot be 0

expression:
    string_literal
    identifier
    identifier = expression
    $number
    expression[expression]
    expression[number]
    identifier(expression1, expression2, ..., expressionN)  // N can be 0
    
    expression || expression
    expression && expression
    expression != expression
    expression == expression
    expression + expression
    
    (expression)
    
    ifelse
    while
    
ifelse:
    if (expression) { block } else { block }
    if (expression) { block } else ifelse                   // This effectively allows else-if 
    
while:
    while (expression) { block }
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
1) The boolean value of a `StringLang` value is `true` 
if and only if the String is not `""` and not `"false"`, else it is `false`.
2) The `StringLang` value of a boolean result of a logic operation is always either `"true"` or `"false"` 

### Evaluation 
```
string_literal              The value of 'string_literal'
identifier                  The current value of the variable with identifier 'identifier'.
identifier = expression     The value of 'expression'. 
                            Side effect: Variable 'identifier' now has that value.
$number                     The value of the 'number'-th (zero-indexed) argument to the program.
expr1[expr2 or number]      The character at position 'expr2' resp. 'number' of the value 
                            that 'expr1' evaluates to.
identifier(expr1, ...)      Function call to built-in function 'identifier' with arguments 'expr1, ...'.
                            Arguments are evaluated before passed to the function.
                            Evaluates to the function's return value.
expr1 <binop> expr2         The value of the corresponding binary operation. 

ifelse                      The value of the then-block if the condition evaluates to a true value,
                            or the value of the else-block if the condition evaluates to a false value.
while                       The value of the last iteration of the block if it gets executed once, else "".
                            Side effects: All iterations might cause side effects
                            
block                       The value of the last expression in the block.
                            Side effects: Evaluates all expressions in the list.

```

#### Exceptions, errors

In the case of an invalid character-access or argument expression, the returned value is always `""`.

### Context (functions and arguments)

To evaluate a `StringLang` program, the interpreter needs an `ast.Context` object.
It contains fields which allow the end-user to supply their custom built-in functions and arguments to the program.
See `cmd/stringlang/main.go` for an example.

There are also some stoppers, mainly a limit for the number of iterations a `while` loop can do 
before being prematurely stopped and one to cap how large the sum of the sizes of all variable-values is allowed to become.

## Why?

It started out as a simple String builder using provided arguments (read: `"my string".replace("$0", args[0])...`) 
for a chat bot, which allowed generating new commands in-chat. 
After switching to a parsed grammar for some simple operators, I figured why not just make it Turing-complete.
What could possibly go wrong allowing users to build their own commands...

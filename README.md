# Ergolas (WIP)

This is just an **e**mbeddable **r**andom **go**lang **la**nguage for **s**cripting. Recently I use Golang very often and sometimes I want to add some extensibility features to my projects using some kind of scripting languages or DSL so I made this mini language for experimenting. There is an included tree walking interpreter but I plan to make it really easy to just parse an expression and get an AST to evaluate with a custom interpreter. 

The syntax is very simple and pretty general and inherits many things from [Lisp](https://en.wikipedia.org/wiki/Lisp_(programming_language)), [REBOL](http://www.rebol.com/)/[Red](https://www.red-lang.org/) and [Julia](https://julialang.org/).

```lua
println "Hello, World!"
```

## Features

- [ ] Parser
    - [x] Function calls without parentheses
    - [x] Property access
    - [x] Quoted forms
    - [x] Binary operators
    - [ ] Unary operators
    - [ ] String templating (for now missing, maybe something can be done just using quoted forms)
    - [ ] Quasi-quotes (Yeah I forgot about them, maybe I will make `"=" <Expr>` be the `unquote` operator)
- [ ] Interpreter
    - [ ] Simple tree walking interpreter
        - [x] Basic operators and arithmetic
        - [x] Basic printing and exiting
        - [x] Basic variable assignment
        - [ ] Lexical scoping
        - [ ] Control flow
        - [ ] Objects and complex values
        - [ ] Dynamic scoping
        - [ ] Hygienic macros
    - [ ] More advanced interpreters...
- [ ] Easily usable as a library
- [ ] Small standard library
- [ ] Interop from and with Go
- [ ] Tooling
    - [ ] Syntax highlighting for common editors
    - [ ] `PKGBUILD` for easy global installation on Arch Linux thorough GitHub releases (mostly for trying this out with GitHub Actions)    

## Usage

To try this out in a REPL (with colors!)

```bash shell
$ go run ./cmd/repl
```

## Reference

### Literals

```perl
# Integer
1

# Decimal
3.14

# Identifier
an-Example_identifier

# String
"an example string"

# List (?) (not implemented)
[1 2 3 4 5] # equivalent to "List 1 2 3 4 5"

# Maps (?) (not implemented)
{ a -> 1, b -> 2, c -> 3 }
```

### Comments

```perl

# This is an inline comment

```

### Functions

Function call don't require parentheses if they 

```perl
# [x] Parses ok, [x] Evals ok
println "a" "b" c" 
```

```perl
# [x] Parses ok, [x] Evals ok
exit 1
```

### Operators

The following binds "a" to 9, arithmetic operators don't have any precedence and are all left associative. There are a only a few right associative operators that for now just are `:=`, `::` even if only `:=` is used for binding variables, `::` will later be used to tell the type of variables.

```perl
# [x] Parses ok, [x] Evals ok
a := 1 + 2 * 3
```

#### Overloading

```perl
# [x] Parses ok, [ ] Evals ok
operator lhs ++ rhs {
    return List.join lhs rhs
}
```

### Misc

Some more examples and ideas for the language syntax and semantics

```perl
# [x] Parses ok, [ ] Evals ok
len := (v.x ^ 2) + (v.y ^ 2)

# [x] Parses ok, [ ] Evals ok
if { a > b } {
    println "True case"
} {
    println "False case" 
}

# [x] Parses ok, [ ] Evals ok
my-list-1 := list 1 2 3 4 5

# [x] Parses ok, [ ] Evals ok
for item my-list-1 {
    printfln "item = {}" item
}
```

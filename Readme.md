[![Build Status](https://drone.io/github.com/davidrjenni/sc/status.png)](https://drone.io/github.com/davidrjenni/sc/latest)
[![GoDoc](https://godoc.org/github.com/davidrjenni/sc?status.svg)](https://godoc.org/github.com/davidrjenni/sc)

# sc - simple compiler

sc is a simple compiler. It compiles a very simple subset of the [Go programming language](http://golang.org).

The language has the types ``bool`` and ``int``. The control structures are ``if``, ``else`` and ``for``.

***

## Language
This is an example of the language:

```
/* Calculates the greatest common factor. */
var a int
var b int
var f int

a = 45
b = 60
for a != b {
    if a < b {
        b = b - a
    }
    if b < a {
        a = a - b
    }
}
f = a
print(f) // output: 15
```

***

### Syntax
The syntax is defined in [Backus-Naur Form](https://en.wikipedia.org/wiki/Backus%E2%80%93Naur_Form).
```
prog    : stmts ;
stmts   : stmt | stmt stmts ;
stmt    : var "=" expr
        | "var" ident type
        | "if" expr "{" stmts "}"
        | "if" expr "{" stmts "}" "else" "{" stmts "}"
        | "for" expr "{" stmts "}"
        | "print" "(" expr ")"
        ;
expr    : "(" expr ")"
        | expr binop expr
        | unop expr
        | var
        | lit
        ;
var     : ident ;
binop   : "+" | "-" | "<" | "<=" | "==" | "!=" | ">=" | ">" | "*" | "/" | "&&" | "||" ;
unop    : "!" | "-" ;
type    : "bool"
        | "int"
        ;
lit     : number
        | "true"
        | "false"
        ;
```

#### Keywords
Following keywords are reserved. They cannot be used as identifiers.

```bool, else, false, for, if, int, print, true, var```

#### Operators and delimiters
Following character sequences act as delimiters.

```{ } ( )```

Following character sequences act as operators.

```+ - / * < <= == != >= > && || !```


#### Comments
There are two forms of comments. They do not nest.

- Line comments terminate at the end of the line. They act like a newline.

```// Characters EOL```

- Block comments can span multiple lines. They act like a space. If they span one or more lines, the act as a newline.

```/* Characters */```

***

### Semantics
- Variables of type ``int`` are initialized with ``0``.
- Variables of type ``bool`` are initialized with ``false``.
- ``if`` and ``for`` statements only accept expressions of type ``bool``.

#### Operators
- Following operators are binary operators.

```* / + - < <= == != >= > && ||```

- Following operators are unary operators.

```- !```

Unary operators have higher precedence than binary operators.
The precedence of the binary operators is the following (in descending order).

1. ``* /``
2. ``+ -``
3. ``< <= == != >= >``
4. ``&&``
5. ``||``

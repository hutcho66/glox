# glox - A lox interpreter written in Go

The lox language was developed by Robert Nystrom for the book [Crafting Interpreters](https://craftinginterpreters.com/).

This is a implementation of the language in go, with a few minor changes.

## Language Specification

Lox is a basic language with some object oriented features.

### Basic Operations

Numbers are implemented as double precision floats but print as integers if it is suitable

```bash
> 5 + 4
9
> 3 / 2
1.500000
```

Strings can be concatenated using the `+` operator

```bash
> "hello " + "world"
hello world
```

All values are truthy except the nil value and boolean false

```bash
> !false
true
> !nil
true
> !0
false
> !""
false
```

Comments can be defined using `//` and must be on one line.

```bash
// this is a comment
var a = 5; // this is another comment
```

### Statement Termination

Lox programs are a sequence of statements. In lox, statements are **not** expressions (even expression statements) and
hence do not have a value. In the reference lox implementation, statements must be terminated with semicolons. However,
unlike the reference lox implementation, in glox, semicolons are optional unless there is multiple statements on a line.

```
> var a = 4 // this is valid
> a = a + 1; print a // also valid
5
> a = 5 b = a // invalid
[line 1] Error at 'b': Improperly terminated statement
```

While statements in lox are not expressions and do not have a value, when using the REPL, if the final statement on a
line is an expression statement (consists solely of an expression), it the value of the expression will be printed.

Note that assignment is an expression but declaration is not.

```
> 5 + 4; 3 - 2;
1
> var a = 5; // declaration is not an expression, so this prints nothing
> a = 3;     // but assignment is, so this prints the new value of a
3
```

### Control Flow and Looping

Lox has `if`-`else` statements which work like any other language. Then and else statements can be singular statements or block statements.

```
> if (6 > 5) print true; else print false;
true

> var x = 5; 
> if (x <= 5) x = x + 1; // assignment in conditional statements is fine

> if (x <= 5) var y = x; // declaration is not allowed in conditional statements
[line 1] Error at 'var': Expect expression. 

> if (x <= 5) { var y = x; } // this is fine, `y` is scoped to the block
```

Lox has C-style `while` and `for` loops. Variables defined in `for` loop initializers are scoped to the loop.
```
> var x = 0;
> while (x < 2) { print x; x = x+1; }
0
1
> x
2

> for (var y = 0; y < 2; y = y+1) print y;
0
1
> y
[line 1] Error at 'y': Undefined variable 'y'
```

## Usage

### Build and run immediately

```bash
# Run the REPL
go run src/cmd/glox.go

# Run a .lox source code file
go run src/cmd/glox.go <path_to_script>
```

### Build and run from local directory

```bash
# Create a binary in bin/
go build -o bin/ src/cmd/glox.go

# Run the binary
./bin/glox
```

### Build and install to GOPATH

```bash
# Build binary and store in GOPATH
go install src/cmd/glox.go

# Run from GOPATH (ensure GOPATH is set in $PATH)
glox
```

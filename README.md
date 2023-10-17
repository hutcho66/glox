# glox - A lox interpreter written in Go

The lox language was developed by Robert Nystrom for the book [Crafting Interpreters](https://craftinginterpreters.com/).

This is a implementation of the language in go, with a few minor changes.

## Language Specification

Lox is a basic language with some object oriented features.

### Basic Operations

Numbers are implemented as double precision floats but print as integers if it is suitable

```js
> 5 + 4
9
> 3 / 2
1.500000
```

Strings can be concatenated using the `+` operator

```js
> "hello " + "world"
hello world
```

All values are truthy except the nil value and boolean false

```js
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

```js
// this is a comment
var a = 5; // this is another comment
```

### Statement Termination

Unlike the reference lox implementation, semicolons are optional unless there is multiple statements on a line.

```js
> var a = 4
> a = a + 1; print a
5
```

When using the REPL, the final statement on a line will be printed if it is an expression (whether terminated with a semicolon or not). Note that assignment is an expression but declaration is not.

```js
> 5 + 4; 3 - 2;
1
> var a = 5; // declaration is not an expression
> a = 3;     // but assignment is
3
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

# glox - A lox interpreter written in Go

![Tests](https://github.com/hutcho66/glox/actions/workflows/go.yml/badge.svg)
[![Go Coverage](https://github.com/hutcho66/glox/wiki/coverage.svg)](https://raw.githack.com/wiki/hutcho66/glox/coverage.html)

The lox language was developed by Robert Nystrom for the book [Crafting Interpreters](https://craftinginterpreters.com/).

This is a implementation of the language in go, with a few additions:

- Optional semicolons - a statement must be terminated either by a semicolon or a newline
- Comma separated sequence expressions
- C-style ternary operator
- Arrays and string-keyed maps
- for..of loops on arrays
- Index notation for accessing arrays, maps and substrings of strings
- break and continue statements within loops
- Lambda expressions using a JavaScript style arrow syntax
- Additonal builtin functions, e.g. `len`, `map`, `filter`, `reduce`
- Classes with static, getter and setter functions


## Table of Contents
  - [Language Specification](#language-specification)
    - [Basic Operations](#basic-operations)
    - [Index notation for strings](#index-notation-for-strings)
    - [Arrays](#arrays)
    - [Maps](#maps)
    - [Statement Termination](#statement-termination)
    - [Control Flow and Looping](#control-flow-and-looping)
    - [Functions](#functions)
    - [Classes](#classes)
  - [Usage](#usage)
    - [Tests](#tests)
    - [Install to GOPATH](#install-to-gopath)
  
## Language Specification

glox is a basic language with some object oriented features.

### Basic Operations

Numbers are implemented as double precision floats, but should print nicely (using go's handy `strconv.FormatFloat` function)

```
> 5 + 4
9
> 3 / 2
1.5
> 5 / 3
1.6666666666666667
```

Strings can be concatenated using the `+` operator. Additionally, a string can be concatenated
with a number or a boolean value

While you cannot concatenate strings with values that are not numbers or booleans, you can get
the string represntation of any value using the builtin `string` function (or print it to the console with the `print` function)

```
> "hello " + "world"
hello world

> "there are " + 9 + " planets in the solar system"
there are 9 planets

> "the last statement is " + false
the last statement is false

> fun a() {}
> "The representation of a is " + a
[line 1] Error at '+': cannot concatenate string with type <fn a>
> "The representation of a is " + string(a)
The representation of a is <fn a>
```

All values are truthy except the nil value and boolean false

```
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

```
// this is a comment
var a = 5 // this is another comment
```

### Index notation for strings

Index notation can be used to retrieve a substring from a string. glox strings are immutable,
so you cannot use index notation to assign to strings.

The `len` builtin function works on strings

```
> var x = "hello"
> x = x[1:5]
"ello"
> x[0] = "E"
[line 1] Error at ']': Can only assign to array elements
> len(x)
4
```

### Arrays

glox supports arrays which are dynamic length and type, and can include any valid glox value. 
Arrays can only be indexed by integers, and are zero-index.

Array assignment is only valid for current length of array, to grow the array you can concatenate two arrays using the `+` operator

Arrays can be accessed (but not assigned) using a slice syntax, `x[1:3]`.

```
> var x = [1, 2, "hello"]
> x[2]
"hello"
> x[2] = 3
3
> x
[1, 2, 3]
> x[3] = 4
[line 1] Error at ']': Index is out of range for array
> x = x + [4]
[1, 2, 3, 4]
> x[0:2]
[1, 2]
> x = x[0:3]
[1, 2, 3]
```

A few builtin functions have been added for arrays
- `len` returns the length of the array
- `map` applies a function to the elements of an array and returns a new array with the results
- `filter` applies a function to the elements of an array and returns a new array with the elements of the original array, if the function returned a truthy value
- `reduce` takes an initial value, an array, and an accumulator function. It applies the function to each element of the array in turn, accumulating the result, beginning with the initial value. The accumulator function must take two parameters: the accumulated value and the element; and return the new accumulated value

```
> var arr = [1,2,3,4]
> len(arr)
4

> map(arr, a => a*2)
[2, 4, 6, 8]

> arr = [1,-2,3,-4]
> filter(arr, a => a > 0)
[1, 3]

> arr = [1,2,3,4,5]
> reduce(0, arr, (acc, el) => acc + el) // should add the elements in the array
15

> var sum = arr => reduce(0, arr, (acc, el) => acc + el)
> sum([1,2,3,4,5])
15
```

### Maps
glox has maps with string keys. Values can be any valid glox value. The default value of maps is `nil`.
```
> var x = {"foo": "bar"}
> x["foo"]
"bar"
> x["bar"] = "foo"
"foo"
> x["goo"]
nil
```

The builtin function `size` gets the number of elements in the map
```
> var x = {"foo": 0, "bar": 1}
> size(x)
2
```

In the case that you want to store `nil` as a valid value in a map, the `hasKey` builtin can be used
to test for the presence of a key
```
> var x = {"foo": nil}
> x["foo"]
nil
> x["bar"]
nil

> hasKey(x, "foo")
true
> hasKey(x, "bar")
false
```

Some care needs to be taken when returning empty maps from lambda expressions
```
> var x = () => {}  // this is a lambda with an empty function body
> x()
nil

> var x = () => { return {} } // this is a lambda that returns an empty map
> x()
<map>
```

### Statement Termination

glox programs are a sequence of statements. In glox, statements are **not** expressions (even expression statements) and
hence do not have a value. In the reference lox implementation, statements must be terminated with semicolons. However,
unlike the reference lox implementation, in glox, semicolons are optional unless there is multiple statements on a line.

```
> var a = 4 // this is valid
> a = a + 1; a // also valid, semicolon is necessary when statements are on the same line
5
> a = 5 b = a // invalid
[line 1] Error at 'b': Improperly terminated statement
```

Statements in glox are not expressions and do not have a value, however when using the REPL, if the final statement on a
line is an expression statement (consists solely of an expression), it the value of the expression will be printed.

Note that assignment is an expression but declaration is not.

```
> 5 + 4; 3 - 2;
1
> var a = 5 // declaration is not an expression, so this prints nothing
> a = 3     // but assignment is, so this prints the new value of a
3
```

### Control Flow and Looping

glox has `if-else` statements which work like any other language. Then and else statements can be singular statements or block statements.

```
> if (6 > 5) print(true); else print(false)
true

> var x = 5 
> if (x <= 5) x = x + 1 // assignment in conditional statements is fine

> if (x <= 5) var y = x // declaration is not allowed in conditional statements
[line 1] Error at 'var': Expect expression. 

> if (x <= 5) { var y = x } // this is fine, `y` is scoped to the block
```

glox also has ternary expressions, which are right associative and at a lower precedence than all other expressions other than lambdas
```
> var x = 1
> x == 1 ? "one" : a == 2 ? "two" : "many"
one
```

glox has C-style `while` and `for` loops. Variables defined in `for` loop initializers are scoped to the loop.
```
> var x = 0
> while (x < 2) { print(x); x = x+1 }
0
1
> x
2

> for (var y = 0; y < 2; y = y+1) print(y)
0
1
> y
[line 1] Error at 'y': Undefined variable 'y'
```

glox supports break and continue statements
```
> var i = 0
> while (i < 10) { if (i == 5) break; i = i + 1 }
> i
5

> for (var i = 0; i < 5; i = i + 1) { if (i == 2) continue; else print(i) }
1
3
4
5
```

Comma separated sequence expressions can be used if there's a need to run more than one expression
in a for loop increment or initializer
```
> var x
> var y
> for ((x=0, y=10); x != y; (x = x+1, y=y-1)) print(x + " " + y)
0 10
1 9
2 8
3 7
4 6
```

Finally, glox supports for..of loops on arrays. This is also useful for looping through maps, 
using the builtin functions `keys` and `values`. Note that maps are unordered.
```
> var arr = [1,2,3,4]
> for (var e of arr) print(e*2)
2
4
6
8

> var mp = {"a": 1, "b": 2, "c": 3}
> for (var key of values(mp)) print(key)
1
3
2

> for (var key of keys(mp)) print(key + ":" + mp[key])
a:1
c:3
b:2
```

### Functions

glox supports both named functions as per the lox reference implementation, as well as lambda expressions, using a JavaScript style arrow syntax. 
Functions are first class objects and can be stored in variables as well as passed as arguments.

The body of a lambda expression can either be a standard block statement, or it can be a single expression,
which will be implicitly returned.

Like JavaScript, single parameter lambdas do not need parentheses, `x => x` is a valid lambda expression. Lambdas with zero or
more than one parameter must have parentheses.

```
> fun hello(first, getLastName) { return "Hello, " + first + " " + getLastName(); }

> fun scott() { return "Scott"; }
> var smith = () => "Smith"  // implicit return
> var jones = () => { return "Jones"; }

> hello("Mark", scott)
Hello, Mark Scott

> hello("Mark", smith)
Hello, Mark Smith

> hello("Mark", jones)
Hello, Mark Jones

> hello("Mark", () => "Taylor")
Hello, Mark Taylor

> hello("Mark", () => return "Green") // this is invalid, lambda bodys can only be a block statement or an expression, return statements are invalid
[line 1] Error at 'return': Expect expression.
```

Closures are fully supported in both named and lambda functions
```
> var adder = a => b => a + b
> var add5 = adder(5)
> add5(6)
11
```

### Classes

glox classes are defined using the `class` keyword. Classes can be constructed using an optional `init` method.
```
> class Foo { 
  init(value) {
    this.bar = value
  }
> var foo = Foo("baz")
> foo.bar()
baz
```

The `static` keyword can be used to define a class method.
```
> class Math {
  static square(r) {
    return r * r
  }
}
> Math.square(5)
25
```

Single inheritance is supported
```
> class A {
  foo() {
    return "bar"
  }
}
> class B > A {
  foo() {
    return "foo" + super.foo()
  }
}
> B().foo()
"foobar"
```

Getters (keyword `get`) and setters (keyword `set`) are supported.
```
> class Foo { 
  set name(first) {
    this.fullName = first + " Smith"
  }

  get formalGreeting {
    return "Good morning, " + this.fullName
  }
}
> var foo = Foo()
> foo.name = "James" // using a setter
> foo.formalGreeting // using a getter
"Good morning, James Smith"
```



## Usage

```bash
# Run the REPL, useful for rapid testing during development
make run

# Build binary
make build

# Run a .lox source code file using the binary
./glox <path_to_script>
```

### Tests

```bash
# Run test suite
make test

# Run test suite and generate coverage HTML report
make coverage
```

### Install to GOPATH

```bash
# Build binary and store in GOPATH
make install

# Run from GOPATH (ensure GOPATH is set in $PATH)
glox
glox <path_to_script>
```

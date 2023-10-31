package interpreter_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/hutcho66/glox/src/pkg/interpreter"
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/parser"
	"github.com/hutcho66/glox/src/pkg/resolver"
	"github.com/hutcho66/glox/src/pkg/scanner"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected any
	}{
		// ignore whitespace
		{"whitespace", "   \t\r 5", 5.0},

		// ignore comments
		{"comment", `5 // comment`, 5.0},
		{"comment - newline", "// comment\n5", 5.0},

		// basic literals
		{"nil literal", "nil", nil},
		{"true literal", "true", true},
		{"false literal", "false", false},
		{"number literal", "5", 5.0},
		{"decimal literal", "55.4", 55.4},
		{"string literal", `"hello world"`, "hello world"},
		{"multiline string literal", `"hello
world"`, "hello\nworld"},

		// array literal
		{"array literal", "[5, true]", interpreter.LoxArray{5.0, true}},

		// map literal
		{"map literal", `{"foo": "bar"}`, interpreter.LoxMap{interpreter.Hash("foo"): interpreter.MapPair{"foo", "bar"}}},
		{"empty map literal", `{}`, interpreter.LoxMap{}},

		// lambda literal
		{"lambda literal", "() => {}", &interpreter.LoxFunction{}},
		{"lambda literal - one param", "a => {}", &interpreter.LoxFunction{}},
		{"lambda literal - multiple params", "(a,b) => {}", &interpreter.LoxFunction{}},

		// basic unary
		{"negation", "-5", -5.0},
		{"not", "!true", false},

		// basic number binary
		{"addition", "4+5", 9.0},
		{"subtration", "4-5", -1.0},
		{"multiplication", "4*5", 20.0},
		{"division", "5/2", 2.5},

		// order of operations
		{"unary lower than sum", "2+-3", -1.0},
		{"factor lower than sum", "2+3*4", 14.0},
		{"grouping", "(2+3)*4", 20.0},

		// basic comparison
		{"greater", "5>5", false},
		{"greater equal", "5>=5", true},
		{"less", "5<5", false},
		{"less equal", "5<=5", true},
		{"equal", "true==true", true},
		{"not equal", "true!=true", false},

		// string binary
		{"string equal", `"hello" == "hello"`, true},
		{"string equal", `"hello" == "world"`, false},
		{"string not equal", `"hello" != "hello"`, false},
		{"string not equal", `"hello" != "world"`, true},
		{"string concatenation", `"hello " + "world"`, "hello world"},
		{"string concatenation with number", `5 + "=x"`, "5=x"},
		{"string concatenation with boolean", `"x: " + true`, "x: true"},

		// array concatenate
		{"array concat", "[5] + [true]", interpreter.LoxArray{5.0, true}},

		// logical expressions
		{"logical and - returns right if left is true", "true and 5.0", 5.0},
		{"logical and - returns left if left is false", "nil and true", nil},
		{"logical or - returns right if left is false", "false or 5.0", 5.0},
		{"logical or - returns left if left is true", "5.0 or false", 5.0},

		// ternary expression
		{"ternary - true", "5 > 4 ? true : false", true},
		{"ternary - false", "5 < 4 ? true : false", false},
		{"ternary - right associative", "true ? 1 : false ? 2 : 3", 1.0},

		// variable declaration and assignment
		{"variable", "var x = 5; x", 5.0},
		{"variable", "var x = 5; x = x + 1", 6.0},

		// sequence expression
		{"sequence", "var x = 5; (x = x + 1, x = x + 1)", 7.0},
		{"empty sequence", "5 == ()", false},

		// indexing
		{"array index get", "var x = [1, 2, 3]; x[1]", 2.0},
		{"array index slice", "var x = [1, 2, 3]; x[1:3]", interpreter.LoxArray{2.0, 3.0}},
		{"array index assign", "var x = [1, 2, 3]; x[1] = 5; x[1]", 5.0},
		{"map index get", `var x = {"foo": "bar"}; x["foo"]`, "bar"},
		{"array index assign", `var x = {"foo": "bar"}; x["foo"] = "baz"; x["foo"]`, "baz"},
		{"string index get", `var x = "hello"; x[1]`, "e"},
		{"string index slice", `var x = "hello"; x[1:5]`, "ello"},

		// conditionals
		{"if - true", "var x = 5; if (x < 6) x = x+1; x", 6.0},
		{"if - false", "var x = 6; if (x < 6) x = x+1; x", 6.0},
		{"if else - true", `var x = 5; if (x < 6) x = x+1
			else x = x-1; x`, 6.0},
		{"if else - false", `var x = 6; if (x < 6) x = x+1
			else x = x-1; x`, 5.0},

		// block scoping
		{"block scope contains outer scope", "var x = 5; {x = 6}\n x", 6.0},
		{"block scope shadows", "var x = 5; {var x = 6}\n x", 5.0},

		// looping
		{"while", "var x = 0; while (x < 5) x = x+1; x", 5.0},
		{"for", "var x = 0; for (var y = 0; y < 5; y = y+1) x = y; x", 4.0},
		{"for - expression initializer", "var x = 0; var y = 0; for (y = 0; y < 5; y = y+1) x = y; x", 4.0},
		{"for - no clauses", "var x = 0; for (;;) break; x", 0.0},
		{"foreach", "var x = 0; var arr = [0,1,2,3,4]; for (var el of arr) x = el; x", 4.0},
		{"foreach - empty array", "var x = -1; var arr = []; for (var el of arr) x = el; x", -1.0},
		{"break", `var x = 0; while (x < 5) {
				x = x + 1
				if (x == 3) break
			}
			x`, 3.0},
		{"continue", `var x = 0; for (var y = 0; y < 5; y = y+1) {
				if (y == 3) continue
				x = x + 1
			}
			x`, 4.0},

		// function declaration
		{"function declaration", "fun x() {}\n x", &interpreter.LoxFunction{}},
		{"lambda declaration", "var x = () => {}; x", &interpreter.LoxFunction{}},

		// function call
		{"function call", "fun x() {}\n x()", nil},
		{"lambda call", "var x = () => {}; x()", nil},

		// return statements
		{"return", "fun x(a,b) { return a+b }\n x(3,5)", 8.0},
		{"lambda implicit return", "var x = (a,b) => a+b; x(3,5)", 8.0},

		// builtins
		{"clock", "clock()", float64(time.Now().UnixMilli() / 1000.0)},

		{"len - array", "len([1,2,3])", 3.0},
		{"len - string", `len("hello")`, 5.0},
		{"size - map", `size({"foo": "bar"})`, 1.0},
		{"hasKey - map", `hasKey({"foo": "bar"}, "foo")`, true},
		{"hasKey - map", `hasKey({"foo": "bar"}, "bar")`, false},
		{"keys - map", `keys({"foo": "bar"})`, interpreter.LoxArray{"foo"}},
		{"values - map", `values({"foo": "bar"})`, interpreter.LoxArray{"bar"}},

		{"map - array", "map([1,2,3], el => el*2)", interpreter.LoxArray{2.0, 4.0, 6.0}},
		{"filter - array", "filter([1,2,3], el => el<3)", interpreter.LoxArray{1.0, 2.0}},
		{"reduce - array", "reduce(1, [1,2,3], (acc,el) => acc*el)", 6.0},

		{"string - nil", `string(nil)`, "nil"},
		{"string - array", `string(["hello", "world"])`, `["hello", "world"]`},
		{"string - map", "string({})", "<map>"},
		{"string - lambda", "string(() => {})", "<lambda>"},
		{"string - named function", "fun a() {}\n string(a)", "<fn a>"},
		{"string - builtin", "string(clock)", "<native fn clock>"},
	}

	for _, c := range cases {
		s := scanner.NewScanner(c.input)
		tokens := s.ScanTokens()

		p := parser.NewParser(tokens)
		statements := p.Parse()

		assert.False(t, lox_error.HadParsingError(), c.name)

		i := interpreter.NewInterpreter()
		r := resolver.NewResolver(i)

		r.Resolve(statements)

		value, ok := i.Interpret(statements)

		assert.True(t, ok, c.name)

		if _, ok := c.expected.(*interpreter.LoxFunction); ok {
			// can't really compare functions so just pass if value is also a function
			assert.IsType(t, c.expected, value, c.name)
		} else {
			assert.Equal(t, c.expected, value, c.name)
		}
	}
}

func TestPrint(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedOutput any
	}{
		{"print num", "print(5)", "5\n"},
		{"print decimal", "print(5.4)", "5.4\n"},
		{"print string", `print("hello")`, "hello\n"},
	}

	for _, c := range cases {
		s := scanner.NewScanner(c.input)
		tokens := s.ScanTokens()

		p := parser.NewParser(tokens)
		statements := p.Parse()

		assert.False(t, lox_error.HadParsingError(), c.name)

		i := interpreter.NewInterpreter()
		r := resolver.NewResolver(i)

		r.Resolve(statements)

		// redirect stdout
		rescueStdout := os.Stdout
		rp, wp, _ := os.Pipe()
		os.Stdout = wp

		i.Interpret(statements)

		wp.Close()
		out, _ := io.ReadAll(rp)
		os.Stdout = rescueStdout

		assert.Equal(t, c.expectedOutput, string(out), c.name)
	}
}

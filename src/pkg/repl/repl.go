package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hutcho66/glox/src/pkg/interpreter"
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/parser"
	"github.com/hutcho66/glox/src/pkg/resolver"
	"github.com/hutcho66/glox/src/pkg/scanner"
)

func RunFile(content string) {
	ipr := interpreter.NewInterpreter()
	run(string(content), ipr, false)

	// If there was an error when parsing, exit before interpreting
	if lox_error.HadParsingError() {
		os.Exit(65)
	}
	if lox_error.HadRuntimeError() {
		os.Exit(70)
	}
}

func RunPrompt() {
	reader := bufio.NewReader(os.Stdin)
	ipr := interpreter.NewInterpreter()
	fmt.Println("Welcome to the glox repl. Press CTRL-Z to exit.")

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		run(line, ipr, true)
		lox_error.ResetError()
	}
}

func run(source string, ipr *interpreter.Interpreter, prompt bool) {
	s := scanner.NewScanner(source)
	toks := s.ScanTokens()

	p := parser.NewParser(toks)
	statements := p.Parse()

	if lox_error.HadParsingError() {
		return
	}

	r := resolver.NewResolver(ipr)
	r.Resolve(statements)

	if lox_error.HadResolutionError() {
		return
	}

	last_expression_value, ok := ipr.Interpret(statements)

	if prompt && ok {
		fmt.Println(interpreter.Stringify(last_expression_value))
	}
}

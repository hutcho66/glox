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
	errors := &lox_error.LoxErrors{}
	ipr := interpreter.NewInterpreter(errors)
	run(string(content), ipr, errors, false)

	// If there was an error when parsing, exit before interpreting
	if errors.HadParsingError() {
		os.Exit(65)
	}
	if errors.HadRuntimeError() {
		os.Exit(70)
	}
}

func RunPrompt() {
	errors := &lox_error.LoxErrors{}

	reader := bufio.NewReader(os.Stdin)
	ipr := interpreter.NewInterpreter(errors)
	fmt.Println("Welcome to the glox repl. Press CTRL-Z to exit.")

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		run(line, ipr, errors, true)
		errors.ResetError()
	}
}

func run(source string, ipr *interpreter.Interpreter, errors *lox_error.LoxErrors, prompt bool) {
	s := scanner.NewScanner(source, errors)
	toks := s.ScanTokens()

	if errors.HadScanningError() {
		return
	}

	p := parser.NewParser(toks, errors)
	statements := p.Parse()

	if errors.HadParsingError() {
		return
	}

	r := resolver.NewResolver(ipr, errors)
	r.Resolve(statements)

	if errors.HadResolutionError() {
		return
	}

	last_expression_value, ok := ipr.Interpret(statements)

	if errors.HadRuntimeError() {
		return
	}

	if prompt && ok {
		fmt.Println(interpreter.Representation(last_expression_value))
	}
}

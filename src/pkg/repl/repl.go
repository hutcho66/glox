package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hutcho66/glox/src/pkg/interpreter"
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/parser"
	"github.com/hutcho66/glox/src/pkg/scanner"
)

func RunFile(content string) {
	interpreter := interpreter.NewInterpreter();
	run(string(content), interpreter);

	// If there was an error when parsing, exit before interpreting
	if lox_error.HadParsingError() {
		os.Exit(65);
	}
	if lox_error.HadRuntimeError() {
		os.Exit(70);
	}
}

func RunPrompt() {
	reader := bufio.NewReader(os.Stdin);
	interpreter := interpreter.NewInterpreter();
	fmt.Println("Welcome to the glox repl. Press CTRL-Z to exit.");

	for {
		fmt.Print("> ");
		line, err := reader.ReadString('\n');
		if err != nil {
			panic(err);
		}
		run(line, interpreter);
		lox_error.ResetError();
	}
}

func run(source string, interpreter *interpreter.Interpreter) {
	s := scanner.NewScanner(source);
	toks := s.ScanTokens();

	p := parser.NewParser(toks);
	statements := p.Parse();

	if lox_error.HadParsingError() {
		return;
	}

	interpreter.Interpret(statements);
}

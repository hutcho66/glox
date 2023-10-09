package interpreter

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/scanner"
)

func RunFile(content string) {
	run(string(content));

	// If there was an error when parsing, exit before interpreting
	if lox_error.HadError() {
		os.Exit(65);
	}
}

func RunPrompt() {
	reader := bufio.NewReader(os.Stdin);
	fmt.Println("Welcome to the glox repl. Press CTRL-Z to exit.");

	for {
		fmt.Print("> ");
		line, err := reader.ReadString('\n');
		if err != nil {
			panic(err);
		}
		run(line);
		lox_error.ResetError();
	}
}

func run(source string) {
	s := scanner.NewScanner(source);
	tokens := s.ScanTokens();

	for _, token := range tokens {
		fmt.Println(token);
	}
}

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hutcho66/glox/src/pkg/interpreter"
)

func main() {
	args := os.Args[1:];
	if len(args) > 1 {
		panic("Usage: glox [args]");
	} else if len(args) == 1 {
		cwd, _ := os.Getwd();
		content, err := os.ReadFile(filepath.Join(cwd, args[0]));
		if err != nil {
			panic(fmt.Sprintf("Invalid path '%s', ensure path is relative to current working directory.", args[0]));
		}
		interpreter.RunFile(string(content));
	} else {
		interpreter.RunPrompt();
	}
}

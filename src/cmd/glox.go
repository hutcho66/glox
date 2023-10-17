package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hutcho66/glox/src/pkg/repl"
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
		repl.RunFile(string(content));
	} else {
		repl.RunPrompt();
	}
}

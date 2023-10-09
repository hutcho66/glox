package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args[1:];
	if len(args) > 1 {
		panic("Usage: glox [args]");
	} else if len(args) == 1 {
		runFile(args[0]);
	}
}

func runFile(path string) {
	cwd, _ := os.Getwd();
	content, err := os.ReadFile(filepath.Join(cwd, path));
	if err != nil {
		panic(fmt.Sprintf("Invalid path '%s', ensure path is relative to current working directory.", path));
	}

	source := string(content);
	fmt.Println(source);
}

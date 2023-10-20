package lox_error

import (
	"errors"
	"fmt"

	"github.com/hutcho66/glox/src/pkg/token"
)

var hadParsingError = false
var hadRuntimeError = false

func ScannerError(line int, message string) {
	hadParsingError = true
	Report(line, "", message)
}

func ParserError(t *token.Token, message string) error {
	hadParsingError = true
	if t.GetType() == token.EOF {
		Report(t.GetLine(), " at end", message)
	} else {
		Report(t.GetLine(), " at '"+t.GetLexeme()+"'", message)
	}

	return errors.New("")
}

func RuntimeError(t *token.Token, message string) error {
	hadRuntimeError = true
	Report(t.GetLine(), " at '"+t.GetLexeme()+"'", message)
	return errors.New("")
}

func Report(line int, where, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}

func HadParsingError() bool {
	return hadParsingError
}

func HadRuntimeError() bool {
	return hadRuntimeError
}

func ResetError() {
	hadParsingError = false
	hadRuntimeError = false
}

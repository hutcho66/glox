package lox_error

import (
	"errors"

	"github.com/fatih/color"
	"github.com/hutcho66/glox/src/pkg/token"
)

type Reporter interface {
	Report(line int, where, message string)
}

type LoxReporter struct{}

func (LoxReporter) Report(line int, where, message string) {
	color.Red("[line %d] Error%s: %s\n", line, where, message)
}

type LoxErrors struct {
	hadScanningError, hadParsingError, hadResolutionError, hadRuntimeError bool
	reporter                                                               Reporter
}

func NewLoxErrors(reporter Reporter) *LoxErrors {
	return &LoxErrors{reporter: reporter}
}

func (l *LoxErrors) ScannerError(line int, message string) {
	l.hadParsingError = true
	l.reporter.Report(line, "", message)
}

func (l *LoxErrors) ParserError(t *token.Token, message string) error {
	l.hadParsingError = true
	if t.Type == token.EOF {
		l.reporter.Report(t.Line, " at end", message)
	} else {
		l.reporter.Report(t.Line, " at '"+t.Lexeme+"'", message)
	}

	return errors.New("")
}

func (l *LoxErrors) ResolutionError(t *token.Token, message string) error {
	l.hadResolutionError = true
	if t.Type == token.EOF {
		l.reporter.Report(t.Line, " at end", message)
	} else {
		l.reporter.Report(t.Line, " at '"+t.Lexeme+"'", message)
	}

	return errors.New("")
}

func (l *LoxErrors) RuntimeError(t *token.Token, message string) error {
	l.hadRuntimeError = true
	l.reporter.Report(t.Line, " at '"+t.Lexeme+"'", message)
	return errors.New("")
}

func (l *LoxErrors) HadScanningError() bool {
	return l.hadParsingError
}

func (l *LoxErrors) HadParsingError() bool {
	return l.hadParsingError
}

func (l *LoxErrors) HadRuntimeError() bool {
	return l.hadRuntimeError
}

func (l *LoxErrors) HadResolutionError() bool {
	return l.hadResolutionError
}

func (l *LoxErrors) ResetError() {
	l.hadScanningError = false
	l.hadParsingError = false
	l.hadRuntimeError = false
	l.hadResolutionError = false
}

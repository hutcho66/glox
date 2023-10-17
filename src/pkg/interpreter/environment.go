package interpreter

import (
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/token"
)

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		values: map[string]any{},
	}
}

func (e *Environment) get(name token.Token) (any, error) {
	if val, ok := e.values[name.GetLexeme()]; ok {
		return val, nil
	}

	return nil, lox_error.RuntimeError(name, "Undefined variable '" + name.GetLexeme() + "'");
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) assign(name token.Token, value any) error {
	if _, ok := e.values[name.GetLexeme()]; ok {
		e.values[name.GetLexeme()] = value;
		return nil;
	}

	return lox_error.RuntimeError(name, "Undefined variable '" + name.GetLexeme() + "'")
}

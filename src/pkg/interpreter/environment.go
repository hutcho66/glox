package interpreter

import (
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/token"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		values:    map[string]any{},
	}
}

func NewEnclosingEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    map[string]any{},
	}
}

func (e *Environment) get(name *token.Token) (any, error) {
	if val, ok := e.values[name.GetLexeme()]; ok {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, lox_error.RuntimeError(name, "Undefined variable '"+name.GetLexeme()+"'")
}

func (e *Environment) getAt(distance int, name string) any {
	return e.ancestor(distance).values[name]
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}

	return environment
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) assign(name *token.Token, value any) error {
	if _, ok := e.values[name.GetLexeme()]; ok {
		e.values[name.GetLexeme()] = value
		return nil
	}

	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return nil
	}

	return lox_error.RuntimeError(name, "Undefined variable '"+name.GetLexeme()+"'")
}

func (e *Environment) assignAt(distance int, name *token.Token, value any) {
	e.ancestor(distance).values[name.GetLexeme()] = value
}

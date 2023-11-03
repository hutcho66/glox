package interpreter

import (
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/token"
)

type LoxInstance struct {
	Class  *LoxClass
	Fields map[string]any
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		Class:  class,
		Fields: make(map[string]any),
	}
}

func (i *LoxInstance) get(name *token.Token) any {
	if field, ok := i.Fields[name.Lexeme]; ok {
		return field
	}

	method := i.Class.findMethod(name.Lexeme)
	if method != nil {
		return method.bind(i)
	}

	panic(lox_error.RuntimeError(name, "Undefined property '"+name.Lexeme+"'."))
}

func (i *LoxInstance) set(name *token.Token, value any) {
	i.Fields[name.Lexeme] = value
}

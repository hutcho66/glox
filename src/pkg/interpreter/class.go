package interpreter

import (
	"github.com/hutcho66/glox/src/pkg/ast"
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/token"
)

type LoxClass struct {
	Name    string
	Methods map[string]*LoxFunction
	Super   *LoxClass
}

func (c LoxClass) Arity() int {
	initializer := c.findMethod("init")
	if initializer == nil {
		return 0
	}

	return initializer.Arity()
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []any) (any, error) {
	instance := NewLoxInstance(c)

	initializer := c.findMethod("init")
	if initializer != nil {
		initializer.bind(instance).Call(interpreter, arguments)
	}

	return instance, nil
}

func (c *LoxClass) findMethod(name string) *LoxFunction {
	if method, ok := c.Methods[name]; ok {
		return method
	}

	if c.Super != nil {
		return c.Super.findMethod(name)
	}

	return nil
}

func (c *LoxClass) get(name *token.Token) any {

	method := c.findMethod(name.Lexeme)

	if method == nil {
		panic(lox_error.RuntimeError(name, "Undefined property '"+name.Lexeme+"'."))
	}

	if method.declaration.Kind != ast.STATIC_METHOD {
		panic(lox_error.RuntimeError(name, "Cannot call non-static method '"+name.Lexeme+"' directly on class."))
	}

	return method.bind(c)

}

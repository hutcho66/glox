package interpreter

import (
	"github.com/hutcho66/glox/src/pkg/ast"
)

type LoxCallable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []any) (any, error)
}

type LoxFunction struct {
	declaration   *ast.FunctionStatement
	closure       *Environment
	isInitializer bool
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []any) (returnValue any, err error) {
	enclosingEnvironment := interpreter.environment
	environment := NewEnclosingEnvironment(f.closure)

	defer func() {
		if val := recover(); val != nil {
			rv, ok := val.(*LoxControl)
			if !ok || rv.controlType != RETURN {
				// repanic
				panic(val)
			}

			if f.isInitializer {
				returnValue = f.closure.getAt(0, "this")
			} else {
				returnValue = rv.value
			}

			interpreter.environment = enclosingEnvironment
			return
		}
	}()

	for i, param := range f.declaration.Params {
		environment.define(param.Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.declaration.Body, environment)

	if f.isInitializer {
		return f.closure.getAt(0, "this"), nil
	}

	return nil, nil
}

func (f LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnclosingEnvironment(f.closure)
	environment.define("this", instance)
	return &LoxFunction{f.declaration, environment, f.isInitializer}
}

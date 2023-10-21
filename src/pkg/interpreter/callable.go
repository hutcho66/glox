package interpreter

import (
	"github.com/hutcho66/glox/src/pkg/ast"
)

type LoxCallable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []any) any
	String() string
}

type LoxFunction struct {
	declaration *ast.FunctionStatement
	closure     *Environment
}

func NewLoxFunction(declaration *ast.FunctionStatement, closure *Environment) *LoxFunction {
	return &LoxFunction{
		closure:     closure,
		declaration: declaration,
	}
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []any) (returnValue any) {
	enclosingEnvironment := interpreter.environment
	environment := NewEnclosingEnvironment(f.closure)

	defer func() {
		if val := recover(); val != nil {
			rv, ok := val.(*LoxControl)
			if !ok || rv.controlType != RETURN {
				// repanic
				panic(val)
			}

			returnValue = rv.value
			interpreter.environment = enclosingEnvironment
			return
		}
	}()

	for i, param := range f.declaration.Parameters() {
		environment.define(param.GetLexeme(), arguments[i])
	}

	interpreter.executeBlock(f.declaration.Body(), environment)

	// if we've reached here, there was no return statement, so implicitly return nil
	return nil
}

func (f LoxFunction) Arity() int {
	return len(f.declaration.Parameters())
}

func (f LoxFunction) String() string {
	if f.declaration.Name() != nil {
		return "<fn " + f.declaration.Name().GetLexeme() + ">"
	} else {
		return "<lambda>"
	}
}

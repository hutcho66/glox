package interpreter

type LoxCallable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []any) (any, error)
}

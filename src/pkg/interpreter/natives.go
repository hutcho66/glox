package interpreter

import (
	"fmt"
	"time"
)

type ClockNative struct{}

func (ClockNative) Arity() int {
	return 0
}

func (ClockNative) Call(interpreter *Interpreter, arguments []any) any {
	return float64(time.Now().UnixMilli() / 1000.0)
}

func (ClockNative) String() string {
	return "<native fn>"
}

func NewClockNative() LoxCallable {
	return ClockNative{}
}

type PrintNative struct{}

func (PrintNative) Arity() int {
	return 1
}

func (PrintNative) Call(interpreter *Interpreter, arguments []any) any {
	fmt.Println(Stringify(arguments[0]))
	return nil
}

func (PrintNative) String() string {
	return "<native fn>"
}

func NewPrintNative() LoxCallable {
	return PrintNative{}
}

type StringNative struct{}

func (StringNative) Arity() int {
	return 1
}

func (StringNative) Call(interpreter *Interpreter, arguments []any) any {
	return Stringify(arguments[0])
}

func (StringNative) String() string {
	return "<native fn>"
}

func NewStringNative() LoxCallable {
	return StringNative{}
}

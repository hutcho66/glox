package interpreter

import (
	"errors"
	"fmt"
	"time"
)

type ClockNative struct{}

func (ClockNative) Arity() int {
	return 0
}

func (ClockNative) Call(interpreter *Interpreter, arguments []any) (any, error) {
	return float64(time.Now().UnixMilli() / 1000.0), nil
}

func (ClockNative) String() string {
	return "<native fn clock>"
}

func NewClockNative() LoxCallable {
	return ClockNative{}
}

type PrintNative struct{}

func (PrintNative) Arity() int {
	return 1
}

func (PrintNative) Call(interpreter *Interpreter, arguments []any) (any, error) {
	fmt.Println(PrintRepresentation(arguments[0]))
	return nil, nil
}

func (PrintNative) String() string {
	return "<native fn print>"
}

func NewPrintNative() LoxCallable {
	return PrintNative{}
}

type StringNative struct{}

func (StringNative) Arity() int {
	return 1
}

func (StringNative) Call(interpreter *Interpreter, arguments []any) (any, error) {
	if s, ok := arguments[0].(string); ok {
		return s, nil
	}
	return Representation(arguments[0]), nil
}

func (StringNative) String() string {
	return "<native fn string>"
}

func NewStringNative() LoxCallable {
	return StringNative{}
}

type LengthNative struct{}

func (LengthNative) Arity() int {
	return 1
}

func (LengthNative) Call(interpreter *Interpreter, arguments []any) (any, error) {
	switch val := arguments[0].(type) {
	case []any:
		return float64(len(val)), nil
	case string:
		return float64(len(val)), nil
	}
	return nil, errors.New("can only call length on arrays or strings")
}

func (LengthNative) String() string {
	return "<native fn len>"
}

func NewLengthNative() LoxCallable {
	return LengthNative{}
}

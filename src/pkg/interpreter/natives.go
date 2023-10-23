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

type MapNative struct{}

func (MapNative) Arity() int {
	return 2
}

func (MapNative) Call(interpreter *Interpreter, arguments []any) (any, error) {
	array, isArray := arguments[0].([]any)
	function, isFunction := arguments[1].(LoxCallable)

	if !isArray {
		return nil, errors.New("first argument of map must be an array")
	}

	if !isFunction || function.Arity() != 1 {
		return nil, errors.New("second argument of map must be an function taking a single parameter")
	}

	results := make([]any, len(array))
	for i, element := range array {
		result, err := function.Call(interpreter, []any{element})
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

func (MapNative) String() string {
	return "<native fn map>"
}

func NewMapNative() LoxCallable {
	return MapNative{}
}

type ReduceNative struct{}

func (ReduceNative) Arity() int {
	return 3
}

func (ReduceNative) Call(interpreter *Interpreter, arguments []any) (any, error) {
	initializer := arguments[0]
	array, isArray := arguments[1].([]any)
	function, isFunction := arguments[2].(LoxCallable)

	if !isArray {
		return nil, errors.New("second argument of reduce must be an array")
	}

	if !isFunction || function.Arity() != 2 {
		return nil, errors.New("third argument of reduce must be an function taking two parameters - the accumulator and the current element")
	}

	accumulator := initializer
	var err error
	for _, element := range array {
		accumulator, err = function.Call(interpreter, []any{accumulator, element})
		if err != nil {
			return nil, err
		}
	}

	return accumulator, nil
}

func (ReduceNative) String() string {
	return "<native fn reduce>"
}

func NewReduceNative() LoxCallable {
	return ReduceNative{}
}

type FilterNative struct{}

func (FilterNative) Arity() int {
	return 2
}

func (FilterNative) Call(interpreter *Interpreter, arguments []any) (any, error) {
	array, isArray := arguments[0].([]any)
	function, isFunction := arguments[1].(LoxCallable)

	if !isArray {
		return nil, errors.New("first argument of map must be an array")
	}

	if !isFunction || function.Arity() != 1 {
		return nil, errors.New("second argument of map must be an function taking a single parameter")
	}

	results := make([]any, 0, len(array))
	for _, element := range array {
		result, err := function.Call(interpreter, []any{element})
		if err != nil {
			return nil, err
		}
		if isTruthy(result) {
			results = append(results, element)
		}
	}

	return results, nil
}

func (FilterNative) String() string {
	return "<native fn filter>"
}

func NewFilterNative() LoxCallable {
	return FilterNative{}
}

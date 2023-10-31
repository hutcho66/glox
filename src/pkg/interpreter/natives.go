package interpreter

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/maps"
)

type LoxNative interface {
	LoxCallable
	Name() string
}

var Natives = []LoxNative{
	&Clock{},
	&Print{},
	&String{},
	&Length{},
	&Map{},
	&Filter{},
	&Reduce{},
	&HasKey{},
	&Size{},
	&Values{},
	&Keys{},
}

type Clock struct{}

func (Clock) Arity() int {
	return 0
}

func (Clock) Call(interpreter *Interpreter, arguments []any) (any, error) {
	return float64(time.Now().UnixMilli() / 1000.0), nil
}

func (Clock) Name() string {
	return "clock"
}

type Print struct{}

func (Print) Arity() int {
	return 1
}

func (Print) Call(interpreter *Interpreter, arguments []any) (any, error) {
	fmt.Println(PrintRepresentation(arguments[0]))
	return nil, nil
}

func (Print) Name() string {
	return "print"
}

type String struct{}

func (String) Arity() int {
	return 1
}

func (String) Call(interpreter *Interpreter, arguments []any) (any, error) {
	if s, ok := arguments[0].(string); ok {
		return s, nil
	}
	return Representation(arguments[0]), nil
}

func (String) Name() string {
	return "string"
}

type Length struct{}

func (Length) Arity() int {
	return 1
}

func (Length) Call(interpreter *Interpreter, arguments []any) (any, error) {
	switch val := arguments[0].(type) {
	case LoxArray:
		return float64(len(val)), nil
	case string:
		return float64(len(val)), nil
	}
	return nil, errors.New("can only call len on arrays or strings")
}

func (Length) Name() string {
	return "len"
}

type Size struct{}

func (Size) Arity() int {
	return 1
}

func (Size) Call(interpreter *Interpreter, arguments []any) (any, error) {
	switch val := arguments[0].(type) {
	case LoxMap:
		return float64(len(val)), nil
	}
	return nil, errors.New("can only call size on maps")
}

func (Size) Name() string {
	return "size"
}

type Map struct{}

func (Map) Arity() int {
	return 2
}

func (Map) Call(interpreter *Interpreter, arguments []any) (any, error) {
	array, isArray := arguments[0].(LoxArray)
	function, isFunction := arguments[1].(LoxCallable)

	if !isArray {
		return nil, errors.New("first argument of map must be an array")
	}

	if !isFunction || function.Arity() != 1 {
		return nil, errors.New("second argument of map must be an function taking a single parameter")
	}

	results := make(LoxArray, len(array))
	for i, element := range array {
		result, err := function.Call(interpreter, []any{element})
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

func (Map) Name() string {
	return "map"
}

type Reduce struct{}

func (Reduce) Arity() int {
	return 3
}

func (Reduce) Call(interpreter *Interpreter, arguments []any) (any, error) {
	initializer := arguments[0]
	array, isArray := arguments[1].(LoxArray)
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

func (Reduce) Name() string {
	return "reduce"
}

type Filter struct{}

func (Filter) Arity() int {
	return 2
}

func (Filter) Call(interpreter *Interpreter, arguments []any) (any, error) {
	array, isArray := arguments[0].(LoxArray)
	function, isFunction := arguments[1].(LoxCallable)

	if !isArray {
		return nil, errors.New("first argument of map must be an array")
	}

	if !isFunction || function.Arity() != 1 {
		return nil, errors.New("second argument of map must be an function taking a single parameter")
	}

	results := make(LoxArray, 0, len(array))
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

func (Filter) Name() string {
	return "filter"
}

type HasKey struct{}

func (HasKey) Arity() int {
	return 2
}

func (HasKey) Call(interpreter *Interpreter, arguments []any) (any, error) {
	m, isMap := arguments[0].(LoxMap)
	key, isString := arguments[1].(string)

	if !isMap {
		return nil, errors.New("first argument of hasKey must be a map")
	}

	if !isString {
		return nil, errors.New("second argument of hasKey must be a string")
	}

	hash := Hash(key)

	_, ok := m[hash]
	return ok, nil
}

func (HasKey) Name() string {
	return "hasKey"
}

type Values struct{}

func (Values) Arity() int {
	return 1
}

func (Values) Call(interpreter *Interpreter, arguments []any) (any, error) {
	m, isMap := arguments[0].(LoxMap)

	if !isMap {
		return nil, errors.New("argument of values must be a map")
	}

	pairs := maps.Values(m)
	values := make(LoxArray, len(pairs))
	for i, pair := range pairs {
		values[i] = pair.Value
	}

	return values, nil
}

func (Values) Name() string {
	return "values"
}

type Keys struct{}

func (Keys) Arity() int {
	return 1
}

func (Keys) Call(interpreter *Interpreter, arguments []any) (any, error) {
	m, isMap := arguments[0].(LoxMap)

	if !isMap {
		return nil, errors.New("argument of keys must be a map")
	}

	pairs := maps.Values(m)
	keys := make(LoxArray, len(pairs))
	for i, pair := range pairs {
		keys[i] = pair.Key
	}

	return keys, nil
}

func (Keys) Name() string {
	return "keys"
}

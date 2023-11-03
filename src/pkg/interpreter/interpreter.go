package interpreter

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/hutcho66/glox/src/pkg/ast"
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/token"
)

type Interpreter struct {
	globals     *Environment
	environment *Environment
	locals      map[ast.Expression]int
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment()

	// add native functions
	for _, fn := range Natives {
		globals.define(fn.Name(), fn)
	}

	return &Interpreter{
		globals:     globals,
		environment: globals,
		locals:      make(map[ast.Expression]int),
	}
}

func (i *Interpreter) Interpret(statements []ast.Statement) (value any, ok bool) {
	defer func() {
		// catch any errors
		if err := recover(); err != nil {
			ok = false
			return
		}
	}()

	for idx, s := range statements {
		if len(statements) >= 1 && idx == len(statements)-1 {
			if es, ok := s.(*ast.ExpressionStatement); ok {
				// if the last statement is an expression statement, return its value
				return i.executeFinalExpressionStatement(es)
			} else {
				// execute as normal if not expression statement
				i.execute(s)
			}
		} else {
			i.execute(s)
		}
	}

	// last statement is not expression statement, so has no return value
	return nil, false
}

func (i *Interpreter) Resolve(expression ast.Expression, depth int) {
	i.locals[expression] = depth
}

func (i *Interpreter) execute(s ast.Statement) (ok bool) {
	s.Accept(i)
	return true
}

func (i *Interpreter) executeFinalExpressionStatement(s *ast.ExpressionStatement) (result any, ok bool) {
	// instead of using ExpressionStatement visitor which returns nil,
	// visit the Expression itself
	return s.Expr.Accept(i), true
}

func (i *Interpreter) executeBlock(s []ast.Statement, environment *Environment) {
	previous := i.environment
	i.environment = environment

	for _, statement := range s {
		if ok := i.execute(statement); !ok {
			// if statement didn't parse, end execution here
			i.environment = previous
			return
		}
	}

	i.environment = previous
}

func (i *Interpreter) VisitBlockStatement(s *ast.BlockStatement) {
	i.executeBlock(s.Statements, NewEnclosingEnvironment(i.environment))
}

func (i *Interpreter) VisitExpressionStatement(s *ast.ExpressionStatement) {
	i.evaluate(s.Expr)
}

func (i *Interpreter) VisitIfStatement(s *ast.IfStatement) {
	conditionResult := i.evaluate(s.Condition)
	if isTruthy(conditionResult) {
		i.execute(s.Consequence)
	} else if s.Alternative != nil {
		i.execute(s.Alternative)
	}
}

func (i *Interpreter) VisitLoopStatement(s *ast.LoopStatement) {
	environment := i.environment

	// catch break statement
	defer func() {
		if val := recover(); val != nil {
			if val != LoxBreak {
				// repanic - not a break statement
				panic(val)
			}

			// this is necessary because break is usually called inside a block
			// and this panic will stop that block exiting properly
			i.environment = environment
		}
	}()

	for isTruthy(i.evaluate(s.Condition)) {
		// this needs to be pushed to a function so that
		// panic-defer works with continue statements
		i.executeLoopBody(s.Body, s.Increment)
	}
}

func (i *Interpreter) VisitForEachStatement(s *ast.ForEachStatement) {
	outerEnvironment := i.environment

	// catch break statement
	defer func() {
		if val := recover(); val != nil {
			if val != LoxBreak {
				// repanic - not a break statement
				panic(val)
			}

			// this is necessary because break is usually called inside a block
			// and this panic will stop that block exiting properly
			i.environment = outerEnvironment
		}
	}()

	// retrieve the array, it must exists in the outer scope
	a := i.evaluate(s.Array)
	array, ok := a.(LoxArray)
	if !ok {
		panic(lox_error.RuntimeError(s.VariableName, "for-of loops are only valid on arrays"))
	}
	if len(array) == 0 {
		return
	}

	// start a new scope and create the loop variable, initialized to first element of array
	i.environment = NewEnclosingEnvironment(i.environment)
	loop_position := 0
	i.environment.define(s.VariableName.Lexeme, array[loop_position])

	// loop through array
	for {
		// execute the loop
		i.executeLoopBody(s.Body, nil)

		// reassign loop variable to next element of array
		loop_position += 1
		if loop_position < len(array) {
			i.environment.assign(s.VariableName, array[loop_position])
		} else {
			// exit loop, all done
			break
		}
	}

	// restore environment
	i.environment = outerEnvironment
}

func (i *Interpreter) executeLoopBody(body ast.Statement, increment ast.Expression) {
	environment := i.environment

	// catch any continue statement - this will only end current loop iteration
	defer func() {
		if val := recover(); val != nil {
			if val != LoxContinue {
				// repanic - not a continue statement
				panic(val)
			}

			// this is necessary because break is usually called inside a block
			// and this panic will stop that block exiting properly
			i.environment = environment

			// ensure increment is run after continue
			if increment != nil {
				i.evaluate(increment)
			}
		}
	}()

	i.execute(body)
	if increment != nil {
		i.evaluate(increment)
	}
}

func (i *Interpreter) VisitVarStatement(s *ast.VarStatement) {
	var value any = nil

	if s.Initializer != nil {
		value = i.evaluate(s.Initializer)
	}

	i.environment.define(s.Name.Lexeme, value)
}

func (i *Interpreter) VisitFunctionStatement(s *ast.FunctionStatement) {
	function := &LoxFunction{declaration: s, closure: i.environment}
	i.environment.define(s.Name.Lexeme, function)
}

func (i *Interpreter) VisitClassStatement(s *ast.ClassStatement) {
	i.environment.define(s.Name.Lexeme, nil)

	methods := map[string]*LoxFunction{}
	for _, method := range s.Methods {
		function := &LoxFunction{method, i.environment, method.Name.Lexeme == "init"}
		methods[method.Name.Lexeme] = function
	}

	class := &LoxClass{Name: s.Name.Lexeme, Methods: methods}
	i.environment.assign(s.Name, class)
}

func (i *Interpreter) VisitReturnStatement(s *ast.ReturnStatement) {
	var value any = nil
	if s.Value != nil {
		value = i.evaluate(s.Value)
	}

	// Using panic to wind back call stack
	panic(LoxReturn(value))
}

func (i *Interpreter) VisitBreakStatement(s *ast.BreakStatement) {
	// Using panic to wind back call stack
	panic(LoxBreak)
}

func (i *Interpreter) VisitContinueStatement(s *ast.ContinueStatement) {
	// Using panic to wind back call stack
	panic(LoxContinue)
}

func (i *Interpreter) evaluate(e ast.Expression) any {
	return e.Accept(i)
}

func (i *Interpreter) VisitTernaryExpression(e *ast.TernaryExpression) any {
	condition := i.evaluate(e.Condition)

	if isTruthy(condition) {
		return i.evaluate(e.Consequence)
	} else {
		return i.evaluate(e.Alternative)
	}
}

func (i *Interpreter) VisitAssignmentExpression(e *ast.AssignmentExpression) any {
	value := i.evaluate(e.Value)

	distance, ok := i.locals[e]
	if ok {
		i.environment.assignAt(distance, e.Name, value)
	} else {
		i.globals.assign(e.Name, value)
	}

	return value
}

func (i *Interpreter) VisitVariableExpression(e *ast.VariableExpression) any {
	return i.lookupVariable(e.Name, e)
}

func (*Interpreter) VisitLiteralExpression(le *ast.LiteralExpression) any {
	return le.Value
}

func (i *Interpreter) VisitGroupedExpression(ge *ast.GroupingExpression) any {
	return i.evaluate(ge.Expr)
}

func (i *Interpreter) VisitSequenceExpression(e *ast.SequenceExpression) any {
	// evaluate all items but only return final one
	var result any
	for _, item := range e.Items {
		result = i.evaluate(item)
	}

	return result
}

func (i *Interpreter) VisitArrayExpression(e *ast.ArrayExpression) any {
	// represent arrays by slices of any
	array := make(LoxArray, len(e.Items))
	for idx, item := range e.Items {
		array[idx] = i.evaluate(item)
	}

	return array
}

func (i *Interpreter) VisitMapExpression(e *ast.MapExpression) any {
	m := make(LoxMap, len(e.Keys))
	for idx := range e.Keys {
		key, isString := i.evaluate(e.Keys[idx]).(string)
		if !isString {
			panic(lox_error.RuntimeError(e.OpeningBrace, "map keys must be strings"))
		}
		hash := Hash(key)
		value := i.evaluate(e.Values[idx])

		m[hash] = MapPair{Key: key, Value: value}
	}

	return m
}

func (i *Interpreter) VisitGetExpression(e *ast.GetExpression) any {
	object := i.evaluate(e.Object)
	if instance, ok := object.(LoxObject); ok {
		property := instance.get(e.Name)

		// if field is a getter method, call it immediately
		if method, ok := property.(*LoxFunction); ok {
			if method.declaration.Kind == ast.GETTER_METHOD {
				value, err := method.Call(i, []any{})
				if err != nil {
					panic(lox_error.RuntimeError(e.Name, err.Error()))
				}

				return value
			}
		}

		// not a getter method, simple return the property
		return property
	}

	panic(lox_error.RuntimeError(e.Name, "Only instances have properties."))
}

func (i *Interpreter) VisitSetExpression(e *ast.SetExpression) any {
	object := i.evaluate(e.Object)
	if instance, ok := object.(*LoxInstance); ok {
		value := i.evaluate(e.Value)

		// check if name refers to a setter
		method := instance.Class.findMethod(e.Name.Lexeme)
		if method != nil && method.declaration.Kind == ast.SETTER_METHOD {
			// bind and call setter method with value
			boundMethod := method.bind(instance)
			_, err := boundMethod.Call(i, []any{value})
			if err != nil {
				panic(lox_error.RuntimeError(e.Name, err.Error()))
			}

		} else {
			instance.set(e.Name, value)
		}

		return value
	}

	panic(lox_error.RuntimeError(e.Name, "Can only set fields on instances."))
}

func (i *Interpreter) VisitThisExpression(e *ast.ThisExpression) any {
	return i.lookupVariable(e.Keyword, e)
}

func (i *Interpreter) arrayIndexExpression(e *ast.IndexExpression) any {
	object := i.evaluate(e.Object)
	leftIndex, leftIsNumber := i.evaluate(e.LeftIndex).(float64)
	var (
		rightIndex    float64
		rightIsNumber bool = false
	)
	if e.RightIndex != nil {
		rightIndex, rightIsNumber = i.evaluate(e.RightIndex).(float64)
	}

	if !leftIsNumber || !isInteger(leftIndex) {
		panic(lox_error.RuntimeError(e.ClosingBracket, "Index must be integer"))
	}

	if rightIsNumber && (!rightIsNumber || !isInteger(rightIndex)) {
		panic(lox_error.RuntimeError(e.ClosingBracket, "Index must be integer"))
	}

	switch val := object.(type) {
	case LoxArray:
		{
			if leftIndex < 0 || int(leftIndex) >= len(val) ||
				(rightIsNumber && (rightIndex < 0 || int(rightIndex) > len(val))) {
				panic(lox_error.RuntimeError(e.ClosingBracket, "Index is out of range"))
			}
			if rightIsNumber && (leftIndex > rightIndex) {
				panic(lox_error.RuntimeError(e.ClosingBracket, "Right index of slice must be greater or equal to left index"))
			}
			if rightIsNumber {
				return val[int(leftIndex):int(rightIndex)]
			} else {
				return val[int(leftIndex)]
			}
		}
	case string:
		{
			if leftIndex < 0 || int(leftIndex) >= len(val) ||
				(rightIsNumber && (rightIndex < 0 || int(rightIndex) > len(val))) {
				panic(lox_error.RuntimeError(e.ClosingBracket, "Index is out of range"))
			}
			if rightIsNumber && (leftIndex > rightIndex) {
				panic(lox_error.RuntimeError(e.ClosingBracket, "Right index of slice must be greater or equal to left index"))
			}
			if rightIsNumber {
				return val[int(leftIndex):int(rightIndex)]
			} else {
				return string(val[int(leftIndex)]) // go will return a byte
			}
		}
	default:
		panic(lox_error.RuntimeError(e.ClosingBracket, "Unreachable"))
	}
}

func (i *Interpreter) mapIndexExpression(e *ast.IndexExpression) any {
	object := i.evaluate(e.Object).(LoxMap)
	key, isString := i.evaluate(e.LeftIndex).(string)

	if e.RightIndex != nil {
		panic(lox_error.RuntimeError(e.ClosingBracket, "Cannot slice maps"))
	}

	if !isString {
		panic(lox_error.RuntimeError(e.ClosingBracket, "Maps can only be indexed with strings"))
	}

	hash := Hash(key)

	return object[hash].Value
}

func (i *Interpreter) VisitIndexExpression(e *ast.IndexExpression) any {
	object := i.evaluate(e.Object)
	switch object.(type) {
	case LoxArray, string:
		return i.arrayIndexExpression(e)
	case LoxMap:
		return i.mapIndexExpression(e)
	}
	panic(lox_error.RuntimeError(e.ClosingBracket, "Can only index arrays, strings and maps"))
}

func (i *Interpreter) arrayIndexedAssignmentExpression(e *ast.IndexedAssignmentExpression) any {
	array, _ := i.evaluate(e.Left.Object).(LoxArray)
	index, isNumber := i.evaluate(e.Left.LeftIndex).(float64)

	// don't need to check for right index as using a slice for assignment is a parser error
	if !isNumber || !isInteger(index) {
		panic(lox_error.RuntimeError(e.Left.ClosingBracket, "Index must be integer"))
	}
	if index < 0 || int(index) >= len(array) {
		panic(lox_error.RuntimeError(e.Left.ClosingBracket, "Index is out of range for array"))
	}

	value := i.evaluate(e.Value)
	array[int(index)] = value
	return value
}

func (i *Interpreter) mapIndexedAssignmentExpression(e *ast.IndexedAssignmentExpression) any {
	m, _ := i.evaluate(e.Left.Object).(LoxMap)
	key, isString := i.evaluate(e.Left.LeftIndex).(string)

	if !isString {
		panic(lox_error.RuntimeError(e.Left.ClosingBracket, "map keys must be strings"))
	}

	hash := Hash(key)
	value := i.evaluate(e.Value)
	m[hash] = MapPair{Key: key, Value: value}
	return value
}

func (i *Interpreter) VisitIndexedAssignmentExpression(e *ast.IndexedAssignmentExpression) any {
	object := i.evaluate(e.Left.Object)
	switch object.(type) {
	case LoxArray:
		return i.arrayIndexedAssignmentExpression(e)
	case LoxMap:
		return i.mapIndexedAssignmentExpression(e)
	}
	panic(lox_error.RuntimeError(e.Left.ClosingBracket, "Can only assign to arrays and maps"))
}

func (i *Interpreter) VisitLogicalExpression(le *ast.LogicalExpression) any {
	left := i.evaluate(le.Left)

	if le.Operator.Type == token.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return i.evaluate(le.Right)
}

func (i *Interpreter) VisitUnaryExpression(ue *ast.UnaryExpression) any {
	right := i.evaluate(ue.Expr)
	operator := ue.Operator

	switch operator.Type {
	case token.BANG:
		return !isTruthy(right)
	case token.MINUS:
		{
			if r, ok := right.(float64); ok {
				return -r
			}
			panic(lox_error.RuntimeError(operator, "Operand must be a number"))
		}
	}

	// Unreachable
	panic(lox_error.RuntimeError(operator, "Unreachable"))
}

func (i *Interpreter) VisitBinaryExpression(be *ast.BinaryExpression) any {
	left := i.evaluate(be.Left)
	right := i.evaluate(be.Right)
	operator := be.Operator

	switch operator.Type {
	// can compare any type with == or != and don't need to type check
	case token.EQUAL_EQUAL:
		return left == right
	case token.BANG_EQUAL:
		return left != right
	// concatenate can be used on any basic types as long as one or more is a string
	case token.PLUS:
		{
			leftNum, leftIsNumber := left.(float64)
			rightNum, rightIsNumber := right.(float64)
			if leftIsNumber && rightIsNumber {
				return leftNum + rightNum
			}

			leftArr, leftIsArray := left.(LoxArray)
			rightArr, rightIsArray := right.(LoxArray)
			if leftIsArray && rightIsArray {
				return append(leftArr, rightArr...)
			}

			leftStr, leftIsString := left.(string)
			rightStr, rightIsString := right.(string)
			if leftIsString && rightIsString {
				return leftStr + rightStr
			} else if leftIsString {
				return concatenate(operator, leftStr, right, false)
			} else if rightIsString {
				return concatenate(operator, rightStr, left, true)
			} else {
				panic(lox_error.RuntimeError(operator, "only valid for two numbers, two strings, two arrays, or one string and a number or boolean"))
			}
		}
	// all other binary operations are only valid on numbers
	default:
		{
			l, lok := left.(float64)
			r, rok := right.(float64)
			if !lok || !rok {
				panic(lox_error.RuntimeError(operator, "only valid for numbers"))
			}
			switch operator.Type {
			case token.MINUS:
				return l - r
			case token.SLASH:
				return l / r
			case token.STAR:
				return l * r
			case token.GREATER:
				return l > r
			case token.GREATER_EQUAL:
				return l >= r
			case token.LESS:
				return l < r
			case token.LESS_EQUAL:
				return l <= r
			}
		}
	}

	// Unreachable
	panic(lox_error.RuntimeError(operator, "Unreachable"))
}

func (i *Interpreter) VisitLambdaExpression(e *ast.LambdaExpression) any {
	return &LoxFunction{declaration: e.Function, closure: i.environment}
}

func (i *Interpreter) VisitCallExpression(e *ast.CallExpression) any {
	callee := i.evaluate(e.Callee)
	argValues := LoxArray{}
	for _, argExpr := range e.Arguments {
		argValues = append(argValues, i.evaluate(argExpr))
	}

	if function, ok := callee.(LoxCallable); ok {
		if len(argValues) != function.Arity() {
			panic(lox_error.RuntimeError(e.ClosingParen, fmt.Sprintf("Expected %d arguments but got %d", function.Arity(), len(argValues))))
		}
		value, err := function.Call(i, argValues)
		if err != nil {
			panic(lox_error.RuntimeError(e.ClosingParen, err.Error()))
		}

		return value
	}
	panic(lox_error.RuntimeError(e.ClosingParen, "Can only call functions and classes"))
}

func (i *Interpreter) lookupVariable(name *token.Token, expression ast.Expression) any {
	if distance, ok := i.locals[expression]; ok {
		// safe to not check for error as the resolver should have done its job...
		return i.environment.getAt(distance, name.Lexeme)
	} else {
		val, err := i.globals.get(name)
		if err == nil {
			return val
		} else {
			panic(err)
		}
	}
}

func concatenate(operator *token.Token, stringValue string, otherValue any, reverse bool) string {
	var other string
	switch otherValue.(type) {
	case float64, bool:
		other = Representation(otherValue)
	default:
		panic(lox_error.RuntimeError(operator, fmt.Sprintf("cannot concatenate string with type %s", Representation(otherValue))))
	}

	if reverse {
		return other + stringValue
	}
	return stringValue + other
}

func isInteger(value float64) bool {
	return value == float64(int(value))
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if v, ok := value.(bool); ok {
		return v
	}
	return true
}

func Representation(v any) string {
	switch v := v.(type) {
	case nil:
		return "nil"
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case LoxArray:
		{
			itemStrings := make([]string, len(v))
			for i, item := range v {
				itemStrings[i] = Representation(item)
			}
			return "[" + strings.Join(itemStrings, ", ") + "]"
		}
	case LoxMap:
		return "<map>"
	case *LoxFunction:
		if v.declaration.Name != nil {
			return "<fn " + v.declaration.Name.Lexeme + ">"
		} else {
			return "<lambda>"
		}
	case *LoxClass:
		return "<class " + v.Name + ">"
	case *LoxInstance:
		return "<object " + v.Class.Name + ">"
	case LoxNative:
		return "<native fn " + v.Name() + ">"
	}

	return "<object>"
}

func PrintRepresentation(v any) string {
	switch v := v.(type) {
	case string:
		return fmt.Sprint(v)
	case nil, bool, float64, LoxArray, LoxCallable, LoxMap:
		return Representation(v)
	}

	return "<object>"
}

func Hash(v string) int {
	h := fnv.New64a()
	h.Write([]byte(v))
	return int(h.Sum64())
}

package interpreter

import (
	"fmt"
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
	globals.define("clock", NewClockNative())
	globals.define("print", NewPrintNative())
	globals.define("string", NewStringNative())
	globals.define("len", NewLengthNative())

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
	return s.Expr().Accept(i), true
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
	i.executeBlock(s.Statements(), NewEnclosingEnvironment(i.environment))
}

func (i *Interpreter) VisitExpressionStatement(s *ast.ExpressionStatement) {
	i.evaluate(s.Expr())
}

func (i *Interpreter) VisitIfStatement(s *ast.IfStatement) {
	conditionResult := i.evaluate(s.Condition())
	if isTruthy(conditionResult) {
		i.execute(s.Consequence())
	} else if s.Alternative() != nil {
		i.execute(s.Alternative())
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

	for isTruthy(i.evaluate(s.Condition())) {
		// this needs to be pushed to a function so that
		// panic-defer works with continue statements
		i.executeLoopBody(s.Body(), s.Increment())
	}
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

	if s.Initializer() != nil {
		value = i.evaluate(s.Initializer())
	}

	i.environment.define(s.Name().GetLexeme(), value)
}

func (i *Interpreter) VisitFunctionStatement(s *ast.FunctionStatement) {
	function := NewLoxFunction(s, i.environment)
	i.environment.define(s.Name().GetLexeme(), function)
}

func (i *Interpreter) VisitReturnStatement(s *ast.ReturnStatement) {
	var value any = nil
	if s.Value() != nil {
		value = i.evaluate(s.Value())
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
	condition := i.evaluate(e.Condition())

	if isTruthy(condition) {
		return i.evaluate(e.Consequence())
	} else {
		return i.evaluate(e.Alternative())
	}
}

func (i *Interpreter) VisitAssignmentExpression(e *ast.AssignmentExpression) any {
	value := i.evaluate(e.Value())

	distance, ok := i.locals[e]
	if ok {
		i.environment.assignAt(distance, e.Name(), value)
	} else {
		i.globals.assign(e.Name(), value)
	}

	return value
}

func (i *Interpreter) VisitVariableExpression(e *ast.VariableExpression) any {
	return i.lookupVariable(e.Name(), e)
}

func (*Interpreter) VisitLiteralExpression(le *ast.LiteralExpression) any {
	return le.Value()
}

func (i *Interpreter) VisitGroupedExpression(ge *ast.GroupingExpression) any {
	return i.evaluate(ge.Expression())
}

func (i *Interpreter) VisitSequenceExpression(e *ast.SequenceExpression) any {
	// evaluate all items but only return final one
	var result any
	for _, item := range e.Items() {
		result = i.evaluate(item)
	}

	return result
}

func (i *Interpreter) VisitArrayExpression(e *ast.ArrayExpression) any {
	// represent arrays by slices of any
	array := make([]any, len(e.Items()))
	for idx, item := range e.Items() {
		array[idx] = i.evaluate(item)
	}

	return array
}

func (i *Interpreter) VisitIndexExpression(e *ast.IndexExpression) any {
	object := i.evaluate(e.Object())
	leftIndex, leftIsNumber := i.evaluate(e.LeftIndex()).(float64)
	var (
		rightIndex    float64
		rightIsNumber bool = false
	)
	if e.RightIndex() != nil {
		rightIndex, rightIsNumber = i.evaluate(e.RightIndex()).(float64)
	}

	if !leftIsNumber || !isInteger(leftIndex) {
		panic(lox_error.RuntimeError(e.ClosingBracket(), "Index must be integer"))
	}

	if rightIsNumber && (!rightIsNumber || !isInteger(rightIndex)) {
		panic(lox_error.RuntimeError(e.ClosingBracket(), "Index must be integer"))
	}

	switch val := object.(type) {
	case []any:
		{
			if leftIndex < 0 || int(leftIndex) >= len(val) ||
				(rightIsNumber && (rightIndex < 0 || int(rightIndex) > len(val))) {
				panic(lox_error.RuntimeError(e.ClosingBracket(), "Index is out of range for array"))
			}
			if rightIsNumber && (leftIndex > rightIndex) {
				panic(lox_error.RuntimeError(e.ClosingBracket(), "Right index of slice must be greater or equal to left index"))
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
				panic(lox_error.RuntimeError(e.ClosingBracket(), "Index is out of range for array"))
			}
			if rightIsNumber && (leftIndex > rightIndex) {
				panic(lox_error.RuntimeError(e.ClosingBracket(), "Right index of slice must be greater or equal to left index"))
			}
			if rightIsNumber {
				return val[int(leftIndex):int(rightIndex)]
			} else {
				return string(val[int(leftIndex)]) // go will return a byte
			}
		}
	}

	panic(lox_error.RuntimeError(e.ClosingBracket(), "Can only index arrays"))
}

func (i *Interpreter) VisitArrayAssignmentExpression(e *ast.ArrayAssignmentExpression) any {
	array, isArray := i.evaluate(e.Left().Object()).([]any)
	index, isNumber := i.evaluate(e.Left().LeftIndex()).(float64)

	// don't need to check for right index as using a slice for assignment is a parser error

	if !isArray {
		panic(lox_error.RuntimeError(e.Left().ClosingBracket(), "Can only assign to array elements"))
	}
	if !isNumber || !isInteger(index) {
		panic(lox_error.RuntimeError(e.Left().ClosingBracket(), "Index must be integer"))
	}
	if index < 0 || int(index) >= len(array) {
		panic(lox_error.RuntimeError(e.Left().ClosingBracket(), "Index is out of range for array"))
	}

	value := i.evaluate(e.Value())
	array[int(index)] = value
	return value
}

func (i *Interpreter) VisitLogicalExpression(le *ast.LogicalExpression) any {
	left := i.evaluate(le.Left())

	if le.Operator().GetType() == token.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return i.evaluate(le.Right())
}

func (i *Interpreter) VisitUnaryExpression(ue *ast.UnaryExpression) any {
	right := i.evaluate(ue.Expression())
	operator := ue.Operator()

	switch operator.GetType() {
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
	left := i.evaluate(be.Left())
	right := i.evaluate(be.Right())
	operator := be.Operator()

	switch operator.GetType() {
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

			leftArr, leftIsArray := left.([]any)
			rightArr, rightIsArray := right.([]any)
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
			switch operator.GetType() {
			case token.MINUS:
				return l - r
			case token.PLUS:
				return l + r
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
	return NewLoxFunction(e.Function(), i.environment)
}

func (i *Interpreter) VisitCallExpression(e *ast.CallExpression) any {
	callee := i.evaluate(e.Callee())
	argValues := []any{}
	for _, argExpr := range e.Arguments() {
		argValues = append(argValues, i.evaluate(argExpr))
	}

	if function, ok := callee.(LoxCallable); ok {
		if len(argValues) != function.Arity() {
			panic(lox_error.RuntimeError(e.ClosingParen(), fmt.Sprintf("Expected %d arguments but got %d", function.Arity(), len(argValues))))
		}
		value, err := function.Call(i, argValues)
		if err != nil {
			panic(lox_error.RuntimeError(e.ClosingParen(), err.Error()))
		}

		return value
	}
	panic(lox_error.RuntimeError(e.ClosingParen(), "Can only call functions and classes"))
}

func (i *Interpreter) lookupVariable(name *token.Token, expression ast.Expression) any {
	if distance, ok := i.locals[expression]; ok {
		// safe to not check for error as the resolver should have done its job...
		return i.environment.getAt(distance, name.GetLexeme())
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
	case []any:
		{
			itemStrings := make([]string, len(v))
			for i, item := range v {
				itemStrings[i] = Representation(item)
			}
			return "[" + strings.Join(itemStrings, ", ") + "]"
		}
	case LoxCallable:
		return v.String()
	}

	return "<object>"
}

func PrintRepresentation(v any) string {
	switch v := v.(type) {
	case string:
		return fmt.Sprint(v)
	case nil, bool, float64, []any, LoxCallable:
		return Representation(v)
	}

	return "<object>"
}

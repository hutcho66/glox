package interpreter

import (
	"fmt"
	"math"

	"github.com/hutcho66/glox/src/pkg/ast"
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/token"
)

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(statements []ast.Statement) (any, bool) {
	for idx, s := range statements {
		if len(statements) >= 1 && idx == len(statements)-1 {
			if es, ok := s.(*ast.ExpressionStatement); ok {
				// if the last statement is an expression statement, return its value
				return i.evaluate(es.Expr()), true
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

func (i *Interpreter) execute(s ast.Statement) (ok bool) {
	// catch any panics and suppress, returning ok=false
	defer func() {
		if err := recover(); err != nil {
			ok = false
			return
		}
	}()

	s.Accept(i)
	return true
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

func (i *Interpreter) VisitPrintStatement(s *ast.PrintStatement) {
	v := i.evaluate(s.Expr())
	fmt.Println(Stringify(v))
}

func (i *Interpreter) VisitVarStatement(s *ast.VarStatement) {
	var value any = nil

	if s.Initializer() != nil {
		value = i.evaluate(s.Initializer())
	}

	i.environment.define(s.Name().GetLexeme(), value)
}

func (i *Interpreter) evaluate(e ast.Expression) any {
	return e.Accept(i)
}

func (i *Interpreter) VisitAssignmentExpression(e *ast.AssignmentExpression) any {
	value := i.evaluate(e.Value())
	i.environment.assign(e.Name(), value)

	return value
}

func (i *Interpreter) VisitVariableExpression(e *ast.VariableExpression) any {
	if value, err := i.environment.get(e.Name()); err == nil {
		return value
	} else {
		panic(err)
	}
}

func (*Interpreter) VisitLiteralExpression(le *ast.LiteralExpression) any {
	return le.Value()
}

func (i *Interpreter) VisitGroupedExpression(ge *ast.GroupingExpression) any {
	return i.evaluate(ge.Expression())
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

	// can compare any type with == or != and don't need to type check
	switch operator.GetType() {
	case token.EQUAL_EQUAL:
		return left == right
	case token.BANG_EQUAL:
		return left != right
	}

	// for non-comparisons, types must match
	switch l := left.(type) {
	case float64:
		{
			if r, ok := right.(float64); ok {
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
			panic(lox_error.RuntimeError(operator, "Operands are of different type"))
		}
	case string:
		{
			if r, ok := right.(string); ok {
				switch operator.GetType() {
				case token.PLUS:
					return l + r
				case token.EQUAL_EQUAL:
					return l == r
				case token.BANG_EQUAL:
					return l != r
				}
			}
			panic(lox_error.RuntimeError(operator, "Operands are of different type"))
		}
	}

	// Unreachable
	panic(lox_error.RuntimeError(operator, "Unreachable"))
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

func Stringify(v any) string {
	switch v := v.(type) {
	case nil:
		return "nil"
	case string:
		return fmt.Sprint(v)
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:
		{
			if math.Mod(v, 1.0) == 0 {
				return fmt.Sprintf("%.0f", v)
			}
			return fmt.Sprintf("%f", v)
		}
	}

	return "unprintable object"
}

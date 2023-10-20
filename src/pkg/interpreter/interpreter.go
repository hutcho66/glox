package interpreter

import (
	"errors"
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
				val, err := i.evaluate(es.Expr())
				if err == nil {
					return val, true
				}
			} else {
				// execute as normal if not expression statement
				i.execute(s)
			}
		} else {
			i.execute(s)
		}
	}
	return nil, false
}

func (i *Interpreter) execute(s ast.Statement) error {
	return s.Accept(i)
}

func (i *Interpreter) executeBlock(s []ast.Statement, environment *Environment) error {
	previous := i.environment
	i.environment = environment

	for _, statement := range s {
		err := i.execute(statement)
		if err != nil {
			i.environment = previous
			return err
		}
	}

	i.environment = previous
	return nil
}

func (i *Interpreter) VisitBlockStatement(s *ast.BlockStatement) error {
	return i.executeBlock(s.Statements(), NewEnclosingEnvironment(i.environment))
}

func (i *Interpreter) VisitExpressionStatement(s *ast.ExpressionStatement) error {
	_, err := i.evaluate(s.Expr())
	return err
}

func (i *Interpreter) VisitIfStatement(s *ast.IfStatement) error {
	conditionResult, err := i.evaluate(s.Condition())
	if err == nil && isTruthy(conditionResult) {
		err = i.execute(s.Consequence())
	} else if s.Alternative() != nil {
		err = i.execute(s.Alternative())
	}

	// err will still be nil if everything succeeded
	return err
}

func (i *Interpreter) VisitPrintStatement(s *ast.PrintStatement) error {
	v, err := i.evaluate(s.Expr())
	if err == nil {
		fmt.Println(Stringify(v))
	}
	return err
}

func (i *Interpreter) VisitVarStatement(s *ast.VarStatement) error {
	var value any = nil
	var err error = nil

	if s.Initializer() != nil {
		value, err = i.evaluate(s.Initializer())
	}

	if err != nil {
		i.environment.define(s.Name().GetLexeme(), value)
	}

	return err
}

func (i *Interpreter) evaluate(e ast.Expression) (any, error) {
	return e.Accept(i)
}

func (i *Interpreter) VisitAssignmentExpression(e *ast.AssignmentExpression) (any, error) {
	value, err := i.evaluate(e.Value())

	if err == nil {
		err = i.environment.assign(e.Name(), value)
	}

	return value, err
}

func (i *Interpreter) VisitVariableExpression(e *ast.VariableExpression) (any, error) {
	return i.environment.get(e.Name())
}

func (*Interpreter) VisitLiteralExpression(le *ast.LiteralExpression) (any, error) {
	return le.Value(), nil
}

func (i *Interpreter) VisitGroupedExpression(ge *ast.GroupingExpression) (any, error) {
	return i.evaluate(ge.Expression())
}

func (i *Interpreter) VisitLogicalExpression(le *ast.LogicalExpression) (any, error) {
	left, err := i.evaluate(le.Left())
	if err != nil {
		return nil, err
	}

	if le.Operator().GetType() == token.OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(le.Right())
}

func (i *Interpreter) VisitUnaryExpression(ue *ast.UnaryExpression) (any, error) {
	right, err := i.evaluate(ue.Expression())
	if err != nil {
		return nil, err
	}

	switch ue.Operator().GetType() {
	case token.BANG:
		return !isTruthy(right), nil
	case token.MINUS:
		{
			if r, ok := right.(float64); ok {
				return -r, nil
			}
			return 0, lox_error.RuntimeError(ue.Operator(), "Operand must be a number")
		}
	}

	// Unreachable
	return 0, errors.New("unreachable error")
}

func (i *Interpreter) VisitBinaryExpression(be *ast.BinaryExpression) (any, error) {
	left, err := i.evaluate(be.Left())
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(be.Right())
	if err != nil {
		return nil, err
	}
	operator := be.Operator()

	// can compare any type with == or != and don't need to type check
	switch operator.GetType() {
	case token.EQUAL_EQUAL:
		return left == right, nil
	case token.BANG_EQUAL:
		return left != right, nil
	}

	// for non-comparisons, types must match
	switch l := left.(type) {
	case float64:
		{
			if r, ok := right.(float64); ok {
				switch operator.GetType() {
				case token.MINUS:
					return l - r, nil
				case token.PLUS:
					return l + r, nil
				case token.SLASH:
					return l / r, nil
				case token.STAR:
					return l * r, nil
				case token.GREATER:
					return l > r, nil
				case token.GREATER_EQUAL:
					return l >= r, nil
				case token.LESS:
					return l < r, nil
				case token.LESS_EQUAL:
					return l <= r, nil
				}
			}
			return nil, lox_error.RuntimeError(operator, "Operands are of different type")
		}
	case string:
		{
			if r, ok := right.(string); ok {
				switch operator.GetType() {
				case token.PLUS:
					return l + r, nil
				case token.EQUAL_EQUAL:
					return l == r, nil
				case token.BANG_EQUAL:
					return l != r, nil
				}
			}
			return nil, lox_error.RuntimeError(operator, "Operands are of different type")
		}
	}

	// Unreachable
	return nil, lox_error.RuntimeError(operator, "Unreachable")
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

package interpreter

import (
	"errors"
	"fmt"
	"math"

	"github.com/hutcho66/glox/src/pkg/ast"
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/token"
)

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(e ast.Expression) {
	result, err := i.evaluate(e);
	if err == nil {
		switch v := result.(type) {
			case nil:     fmt.Println("nil");
			case string: 	fmt.Println(v);
			case bool:    fmt.Printf("%t\n", v);
			case float64: {
				if math.Mod(v, 1.0) == 0 {
					fmt.Printf("%.0f\n", v);
				} else {
					fmt.Printf("%f\n", v);
				}
			}
		}
	}
}

func (i *Interpreter) evaluate(e ast.Expression) (any, error) {
	return e.Accept(i)
}

func (*Interpreter) VisitLiteralExpression(le *ast.LiteralExpression) (any, error) {
	return le.Value(), nil
}

func (i *Interpreter) VisitGroupedExpression(ge *ast.GroupingExpression) (any, error) {
	return i.evaluate(ge.Expression())
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

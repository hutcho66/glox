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

func (i *Interpreter) Interpret(statements []ast.Statement) {
	for _, s := range statements {
		i.execute(s);
	}
}

func (i *Interpreter) execute(s ast.Statement) error {
	return s.Accept(i);
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
	i.Evaluate(s.Expr())
	return nil
}

func (i *Interpreter) VisitPrintStatement(s *ast.PrintStatement) error {
	v, err := i.Evaluate(s.Expr())
	if err == nil {
		fmt.Println(Stringify(v))
	}
	return nil
}

func (i *Interpreter) VisitVarStatement(s *ast.VarStatement) error {
	var value any = nil;
	var err error = nil;
	if s.Initializer() != nil {
		value, err = i.Evaluate(s.Initializer())
		if err != nil {
			return err
		}
	}

	i.environment.define(s.Name().GetLexeme(), value);
	return nil
}

func (i *Interpreter) Evaluate(e ast.Expression) (any, error) {
	return e.Accept(i)
}

func (i *Interpreter) VisitAssignmentExpression(e *ast.AssignmentExpression) (any, error) {
	value, err := i.Evaluate(e.Value());
	if err != nil {
		return nil, err
	}

	i.environment.assign(e.Name(), value);
	return value, nil
}

func (i *Interpreter) VisitVariableExpression(e *ast.VariableExpression) (any, error) {
	return i.environment.get(e.Name())
}

func (*Interpreter) VisitLiteralExpression(le *ast.LiteralExpression) (any, error) {
	return le.Value(), nil
}

func (i *Interpreter) VisitGroupedExpression(ge *ast.GroupingExpression) (any, error) {
	return i.Evaluate(ge.Expression())
}

func (i *Interpreter) VisitUnaryExpression(ue *ast.UnaryExpression) (any, error) {
	right, err := i.Evaluate(ue.Expression())
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
	left, err := i.Evaluate(be.Left())
	if err != nil {
		return nil, err
	}
	right, err := i.Evaluate(be.Right())
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
		case nil:     return "nil";
		case string: 	return fmt.Sprint(v);
		case bool:    return fmt.Sprintf("%t", v);
		case float64: {
			if math.Mod(v, 1.0) == 0 {
				return fmt.Sprintf("%.0f", v);
			}
			return fmt.Sprintf("%f", v);
		}
	}

	return "unprintable object"
}

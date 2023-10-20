package ast

import (
	"fmt"
	"strings"

	"github.com/hutcho66/glox/src/pkg/token"
)

type Expression interface {
	expression() bool
	String() string
	Accept(ExpressionVisitor) any
}

type ExpressionVisitor interface {
	VisitBinaryExpression(*BinaryExpression) any
	VisitLogicalExpression(*LogicalExpression) any
	VisitGroupedExpression(*GroupingExpression) any
	VisitUnaryExpression(*UnaryExpression) any
	VisitLiteralExpression(*LiteralExpression) any
	VisitVariableExpression(*VariableExpression) any
	VisitAssignmentExpression(*AssignmentExpression) any
	VisitCallExpression(*CallExpression) any
}

type BinaryExpression struct {
	left, right Expression
	operator    token.Token
}

func NewBinaryExpression(left Expression, operator token.Token, right Expression) Expression {
	return &BinaryExpression{
		left:     left,
		right:    right,
		operator: operator,
	}
}

// BinaryExpression implements Expression
func (BinaryExpression) expression() bool {
	return true
}

func (b BinaryExpression) String() string {
	return fmt.Sprintf("%s %s %s", b.left, b.operator.GetLexeme(), b.right)
}

func (b *BinaryExpression) Accept(v ExpressionVisitor) any {
	return v.VisitBinaryExpression(b)
}

func (b BinaryExpression) Left() Expression {
	return b.left
}

func (b BinaryExpression) Right() Expression {
	return b.right
}

func (b BinaryExpression) Operator() token.Token {
	return b.operator
}

type LogicalExpression struct {
	left, right Expression
	operator    token.Token
}

func NewLogicalExpression(left Expression, operator token.Token, right Expression) Expression {
	return &LogicalExpression{
		left:     left,
		right:    right,
		operator: operator,
	}
}

// LogicalExpression implements Expression
func (LogicalExpression) expression() bool {
	return true
}

func (b LogicalExpression) String() string {
	return fmt.Sprintf("%s %s %s", b.left, b.operator.GetLexeme(), b.right)
}

func (b *LogicalExpression) Accept(v ExpressionVisitor) any {
	return v.VisitLogicalExpression(b)
}

func (b LogicalExpression) Left() Expression {
	return b.left
}

func (b LogicalExpression) Right() Expression {
	return b.right
}

func (b LogicalExpression) Operator() token.Token {
	return b.operator
}

type UnaryExpression struct {
	expr     Expression
	operator token.Token
}

func NewUnaryExpression(operator token.Token, expr Expression) Expression {
	return &UnaryExpression{
		expr:     expr,
		operator: operator,
	}
}

// UnaryExpression implements Expression
func (UnaryExpression) expression() bool {
	return true
}

func (u UnaryExpression) String() string {
	return fmt.Sprintf("%s%s", u.operator.GetLexeme(), u.expr)
}

func (u *UnaryExpression) Accept(v ExpressionVisitor) any {
	return v.VisitUnaryExpression(u)
}

func (u UnaryExpression) Expression() Expression {
	return u.expr
}

func (u UnaryExpression) Operator() token.Token {
	return u.operator
}

type GroupingExpression struct {
	expr Expression
}

func NewGroupingExpression(expr Expression) Expression {
	return &GroupingExpression{
		expr: expr,
	}
}

// GroupingExpression implements Expression
func (GroupingExpression) expression() bool {
	return true
}

func (g GroupingExpression) String() string {
	return fmt.Sprintf("(%s)", g.expr)
}

func (g *GroupingExpression) Accept(v ExpressionVisitor) any {
	return v.VisitGroupedExpression(g)
}

func (g *GroupingExpression) Expression() Expression {
	return g.expr
}

type LiteralExpression struct {
	value any
}

func NewLiteralExpression(value any) Expression {
	return &LiteralExpression{
		value: value,
	}
}

// LiteralExpression implements Expression
func (LiteralExpression) expression() bool {
	return true
}

func (l LiteralExpression) String() string {
	switch v := l.value.(type) {
	case float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return "nil"
	case string:
		return v
	default:
		return fmt.Sprintf("%s", v)
	}
}

func (l *LiteralExpression) Accept(v ExpressionVisitor) any {
	return v.VisitLiteralExpression(l)
}

func (l LiteralExpression) Value() any {
	return l.value
}

type VariableExpression struct {
	name token.Token
}

func NewVariableExpression(name token.Token) Expression {
	return &VariableExpression{
		name: name,
	}
}

// VariableExpression implements Expression
func (VariableExpression) expression() bool {
	return true
}

func (e VariableExpression) String() string {
	return e.name.GetLexeme()
}

func (e *VariableExpression) Accept(v ExpressionVisitor) any {
	return v.VisitVariableExpression(e)
}

func (e VariableExpression) Name() token.Token {
	return e.name
}

type AssignmentExpression struct {
	name  token.Token
	value Expression
}

func NewAssignmentExpression(name token.Token, value Expression) Expression {
	return &AssignmentExpression{
		name:  name,
		value: value,
	}
}

// AssignmentExpression implements Expression
func (AssignmentExpression) expression() bool {
	return true
}

func (e AssignmentExpression) String() string {
	return e.name.GetLexeme() + " = " + e.value.String()
}

func (e *AssignmentExpression) Accept(v ExpressionVisitor) any {
	return v.VisitAssignmentExpression(e)
}

func (e AssignmentExpression) Name() token.Token {
	return e.name
}

func (e AssignmentExpression) Value() Expression {
	return e.value
}

type CallExpression struct {
	callee       Expression
	arguments    []Expression
	closingParen token.Token
}

func NewCallExpression(callee Expression, arguments []Expression, closingParen token.Token) Expression {
	return &CallExpression{
		callee:       callee,
		arguments:    arguments,
		closingParen: closingParen,
	}
}

// CallExpression implements Expression
func (CallExpression) expression() bool {
	return true
}

func (e CallExpression) String() string {
	args := []string{}
	for _, arg := range e.arguments {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", e.callee.String(), strings.Join(args, ", "))
}

func (e *CallExpression) Accept(v ExpressionVisitor) any {
	return v.VisitCallExpression(e)
}

func (e CallExpression) ClosingParen() token.Token {
	return e.closingParen
}

func (e CallExpression) Callee() Expression {
	return e.callee
}

func (e CallExpression) Arguments() []Expression {
	return e.arguments
}

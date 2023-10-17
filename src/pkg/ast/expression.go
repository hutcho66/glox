package ast

import (
	"fmt"

	"github.com/hutcho66/glox/src/pkg/token"
)

type Expression interface {
	expression() bool
	String() string
	Accept(Visitor) (any, error)
}

type Visitor interface {
	VisitBinaryExpression(*BinaryExpression) (any, error)
	VisitGroupedExpression(*GroupingExpression) (any, error)
	VisitUnaryExpression(*UnaryExpression) (any, error)
	VisitLiteralExpression(*LiteralExpression) (any, error)
}

type BinaryExpression struct {
	left, right Expression
	operator token.Token
}

func NewBinaryExpression(left Expression, operator token.Token, right Expression) Expression {
	return &BinaryExpression{
		left: left,
		right: right,
		operator: operator,
	};
}

// BinaryExpression implements Expression
func (BinaryExpression) expression() bool {
	return true;
}

func (b BinaryExpression) String() string {
	return fmt.Sprintf("%s %s %s", b.left, b.operator.GetLexeme(), b.right);
}

func (b *BinaryExpression) Accept(v Visitor) (any, error) {
	return v.VisitBinaryExpression(b);
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

type UnaryExpression struct {
	expr Expression
	operator token.Token
}

func NewUnaryExpression(operator token.Token, expr Expression) Expression {
	return &UnaryExpression{
		expr: expr,
		operator: operator,
	};
}

// UnaryExpression implements Expression
func (UnaryExpression) expression() bool {
	return true;
}

func (u UnaryExpression) String() string {
	return fmt.Sprintf("%s%s", u.operator.GetLexeme(), u.expr);
}

func (u *UnaryExpression) Accept(v Visitor) (any, error) {
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
	};
}

// GroupingExpression implements Expression
func (GroupingExpression) expression() bool {
	return true;
}

func (g GroupingExpression) String() string {
	return fmt.Sprintf("(%s)", g.expr);
}

func (g *GroupingExpression) Accept(v Visitor) (any, error) {
	return v.VisitGroupedExpression(g)
}

func (g *GroupingExpression) Expression() Expression {
	return g.expr;
}

type LiteralExpression struct {
	value any
}

func NewLiteralExpression(value any) Expression {
	return &LiteralExpression{
		value: value,
	};
}

// LiteralExpression implements Expression
func (LiteralExpression) expression() bool {
	return true;
}

func (l LiteralExpression) String() string {
	switch v := l.value.(type) {
		case float64: return fmt.Sprintf("%.2f", v);
		case bool:    return fmt.Sprintf("%t", v);
		case nil:     return "nil";
		case string: 	return v;
		default: return fmt.Sprintf("%s", v);
	}
}

func (l *LiteralExpression) Accept(v Visitor) (any, error) {
	return v.VisitLiteralExpression(l);
}

func (l LiteralExpression) Value() any {
	return l.value;
}



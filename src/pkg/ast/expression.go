package ast

import (
	"github.com/hutcho66/glox/src/pkg/token"
)

type Expression interface {
	Accept(ExpressionVisitor) any
}

type ExpressionVisitor interface {
	VisitBinaryExpression(*BinaryExpression) any
	VisitTernaryExpression(*TernaryExpression) any
	VisitLogicalExpression(*LogicalExpression) any
	VisitGroupedExpression(*GroupingExpression) any
	VisitUnaryExpression(*UnaryExpression) any
	VisitLiteralExpression(*LiteralExpression) any
	VisitVariableExpression(*VariableExpression) any
	VisitAssignmentExpression(*AssignmentExpression) any
	VisitCallExpression(*CallExpression) any
	VisitLambdaExpression(*LambdaExpression) any
	VisitSequenceExpression(*SequenceExpression) any
	VisitArrayExpression(*ArrayExpression) any
	VisitMapExpression(*MapExpression) any
	VisitIndexExpression(*IndexExpression) any
	VisitIndexedAssignmentExpression(*IndexedAssignmentExpression) any
}

type BinaryExpression struct {
	Left, Right Expression
	Operator    *token.Token
}

func (b *BinaryExpression) Accept(v ExpressionVisitor) any {
	return v.VisitBinaryExpression(b)
}

type TernaryExpression struct {
	Condition, Consequence, Alternative Expression
	Operator                            *token.Token
}

func (e *TernaryExpression) Accept(v ExpressionVisitor) any {
	return v.VisitTernaryExpression(e)
}

type LogicalExpression struct {
	Left, Right Expression
	Operator    *token.Token
}

func (b *LogicalExpression) Accept(v ExpressionVisitor) any {
	return v.VisitLogicalExpression(b)
}

type UnaryExpression struct {
	Expr     Expression
	Operator *token.Token
}

func (u *UnaryExpression) Accept(v ExpressionVisitor) any {
	return v.VisitUnaryExpression(u)
}

type GroupingExpression struct {
	Expr Expression
}

func (g *GroupingExpression) Accept(v ExpressionVisitor) any {
	return v.VisitGroupedExpression(g)
}

type LiteralExpression struct {
	Value any
}

func (l *LiteralExpression) Accept(v ExpressionVisitor) any {
	return v.VisitLiteralExpression(l)
}

type VariableExpression struct {
	Name *token.Token
}

func (e *VariableExpression) Accept(v ExpressionVisitor) any {
	return v.VisitVariableExpression(e)
}

type AssignmentExpression struct {
	Name  *token.Token
	Value Expression
}

func (e *AssignmentExpression) Accept(v ExpressionVisitor) any {
	return v.VisitAssignmentExpression(e)
}

type IndexedAssignmentExpression struct {
	Left  *IndexExpression
	Value Expression
}

func (e *IndexedAssignmentExpression) Accept(v ExpressionVisitor) any {
	return v.VisitIndexedAssignmentExpression(e)
}

type CallExpression struct {
	Callee       Expression
	Arguments    []Expression
	ClosingParen *token.Token
}

func (e *CallExpression) Accept(v ExpressionVisitor) any {
	return v.VisitCallExpression(e)
}

type LambdaExpression struct {
	Operator *token.Token
	Function *FunctionStatement
}

func (e *LambdaExpression) Accept(v ExpressionVisitor) any {
	return v.VisitLambdaExpression(e)
}

type SequenceExpression struct {
	Items []Expression
}

func (e *SequenceExpression) Accept(v ExpressionVisitor) any {
	return v.VisitSequenceExpression(e)
}

type ArrayExpression struct {
	Items []Expression
}

func (e *ArrayExpression) Accept(v ExpressionVisitor) any {
	return v.VisitArrayExpression(e)
}

type MapExpression struct {
	OpeningBrace *token.Token
	Keys         []Expression
	Values       []Expression
}

func (e *MapExpression) Accept(v ExpressionVisitor) any {
	return v.VisitMapExpression(e)
}

type IndexExpression struct {
	Object         Expression
	LeftIndex      Expression
	RightIndex     Expression
	ClosingBracket *token.Token
}

func (e *IndexExpression) Accept(v ExpressionVisitor) any {
	return v.VisitIndexExpression(e)
}

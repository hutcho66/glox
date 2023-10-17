package ast

import "github.com/hutcho66/glox/src/pkg/token"

type Statement interface {
	statement() bool
	String() string
	Accept(StatementVisitor) error
}

type StatementVisitor interface {
	VisitExpressionStatement(*ExpressionStatement) error
	VisitPrintStatement(*PrintStatement) error
	VisitVarStatement(*VarStatement) error
}

type ExpressionStatement struct {
	expr Expression
}

func NewExpressionStatement(expr Expression) *ExpressionStatement {
	return &ExpressionStatement{
		expr: expr,
	}
}

func (ExpressionStatement) statement() bool {
	return true
}

func (s ExpressionStatement) String() string {
	return s.expr.String() + ";"
}

func (s ExpressionStatement) Expr() Expression {
	return s.expr
}

func (s *ExpressionStatement) Accept(v StatementVisitor) error {
	return v.VisitExpressionStatement(s);
}

type PrintStatement struct {
	expr Expression
}

func NewPrintStatement(expr Expression) *PrintStatement {
	return &PrintStatement{
		expr: expr,
	}
}

func (PrintStatement) statement() bool {
	return true
}

func (s PrintStatement) String() string {
	return "print " + s.expr.String() + ";"
}

func (s PrintStatement) Expr() Expression {
	return s.expr
}

func (s *PrintStatement) Accept(v StatementVisitor) error {
	return v.VisitPrintStatement(s);
}

type VarStatement struct {
	name token.Token
	initializer Expression
}

func NewVarStatement(name token.Token, initializer Expression) *VarStatement {
	return &VarStatement{
		name: name,
		initializer: initializer,
	}
}

func (VarStatement) statement() bool {
	return true
}

func (s VarStatement) String() string {
	return "var " + s.name.GetLexeme() + " = " + s.initializer.String() + ";"
}

func (s VarStatement) Initializer() Expression {
	return s.initializer
}

func (s VarStatement) Name() token.Token {
	return s.name
}

func (s *VarStatement) Accept(v StatementVisitor) error {
	return v.VisitVarStatement(s);
}

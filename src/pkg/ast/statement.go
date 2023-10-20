package ast

import (
	"fmt"
	"strings"

	"github.com/hutcho66/glox/src/pkg/token"
)

type Statement interface {
	statement() bool
	String() string
	Accept(StatementVisitor) error
}

type StatementVisitor interface {
	VisitExpressionStatement(*ExpressionStatement) error
	VisitPrintStatement(*PrintStatement) error
	VisitVarStatement(*VarStatement) error
	VisitBlockStatement(*BlockStatement) error
	VisitIfStatement(*IfStatement) error
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
	return v.VisitExpressionStatement(s)
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
	return v.VisitPrintStatement(s)
}

type VarStatement struct {
	name        token.Token
	initializer Expression
}

func NewVarStatement(name token.Token, initializer Expression) *VarStatement {
	return &VarStatement{
		name:        name,
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
	return v.VisitVarStatement(s)
}

type BlockStatement struct {
	statements []Statement
}

func NewBlockStatement(statements []Statement) *BlockStatement {
	return &BlockStatement{
		statements: statements,
	}
}

func (BlockStatement) statement() bool {
	return true
}

func (s BlockStatement) String() string {
	buf := []string{}
	buf = append(buf, "{")
	for _, statement := range s.statements {
		buf = append(buf, "\t"+statement.String())
	}
	buf = append(buf, "}")
	return strings.Join(buf, "\n")
}

func (s BlockStatement) Statements() []Statement {
	return s.statements
}

func (s *BlockStatement) Accept(v StatementVisitor) error {
	return v.VisitBlockStatement(s)
}

type IfStatement struct {
	condition                Expression
	consequence, alternative Statement
}

func NewIfStatement(condition Expression, consequence, alternative Statement) *IfStatement {
	return &IfStatement{
		condition:   condition,
		consequence: consequence,
		alternative: alternative,
	}
}

func (IfStatement) statement() bool {
	return true
}

func (s IfStatement) String() string {
	return fmt.Sprintf("if (%s) %s else %s", s.condition.String(), s.consequence.String(), s.alternative.String())
}

func (s IfStatement) Condition() Expression {
	return s.condition
}

func (s IfStatement) Consequence() Statement {
	return s.consequence
}

func (s IfStatement) Alternative() Statement {
	return s.alternative
}

func (s *IfStatement) Accept(v StatementVisitor) error {
	return v.VisitIfStatement(s)
}

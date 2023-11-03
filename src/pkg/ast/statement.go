package ast

import (
	"github.com/hutcho66/glox/src/pkg/token"
)

type Statement interface {
	Accept(StatementVisitor)
}

type StatementVisitor interface {
	VisitExpressionStatement(*ExpressionStatement)
	VisitVarStatement(*VarStatement)
	VisitBlockStatement(*BlockStatement)
	VisitIfStatement(*IfStatement)
	VisitLoopStatement(*LoopStatement)
	VisitForEachStatement(*ForEachStatement)
	VisitFunctionStatement(*FunctionStatement)
	VisitReturnStatement(*ReturnStatement)
	VisitBreakStatement(*BreakStatement)
	VisitContinueStatement(*ContinueStatement)
	VisitClassStatement(*ClassStatement)
}

type ExpressionStatement struct {
	Expr Expression
}

func (s *ExpressionStatement) Accept(v StatementVisitor) {
	v.VisitExpressionStatement(s)
}

type VarStatement struct {
	Name        *token.Token
	Initializer Expression
}

func (s *VarStatement) Accept(v StatementVisitor) {
	v.VisitVarStatement(s)
}

type BlockStatement struct {
	Statements []Statement
}

func (s *BlockStatement) Accept(v StatementVisitor) {
	v.VisitBlockStatement(s)
}

type IfStatement struct {
	Condition                Expression
	Consequence, Alternative Statement
}

func (s *IfStatement) Accept(v StatementVisitor) {
	v.VisitIfStatement(s)
}

type LoopStatement struct {
	Condition Expression
	Body      Statement
	Increment Expression
}

func (s *LoopStatement) Accept(v StatementVisitor) {
	v.VisitLoopStatement(s)
}

type ForEachStatement struct {
	VariableName *token.Token
	Array        Expression
	Body         Statement
}

func (s *ForEachStatement) Accept(v StatementVisitor) {
	v.VisitForEachStatement(s)
}

type MethodType int

const (
	NOT_METHOD = iota
	NORMAL_METHOD
	STATIC_METHOD
	GETTER_METHOD
	SETTER_METHOD
)

type FunctionStatement struct {
	Name   *token.Token
	Params []*token.Token
	Body   []Statement
	Kind   MethodType
}

func (s *FunctionStatement) Accept(v StatementVisitor) {
	v.VisitFunctionStatement(s)
}

type ReturnStatement struct {
	Keyword *token.Token
	Value   Expression
}

func (s *ReturnStatement) Accept(v StatementVisitor) {
	v.VisitReturnStatement(s)
}

type BreakStatement struct {
	Keyword *token.Token
}

func (s *BreakStatement) Accept(v StatementVisitor) {
	v.VisitBreakStatement(s)
}

type ContinueStatement struct {
	Keyword *token.Token
}

func (s *ContinueStatement) Accept(v StatementVisitor) {
	v.VisitContinueStatement(s)
}

type ClassStatement struct {
	Name       *token.Token
	Methods    []*FunctionStatement
	Superclass *VariableExpression
}

func (s *ClassStatement) Accept(v StatementVisitor) {
	v.VisitClassStatement(s)
}

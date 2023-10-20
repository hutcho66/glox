package ast

import (
	"fmt"
	"strings"

	"github.com/hutcho66/glox/src/pkg/token"
)

type Statement interface {
	statement() bool
	String() string
	Accept(StatementVisitor)
}

type StatementVisitor interface {
	VisitExpressionStatement(*ExpressionStatement)
	VisitVarStatement(*VarStatement)
	VisitBlockStatement(*BlockStatement)
	VisitIfStatement(*IfStatement)
	VisitWhileStatement(*WhileStatement)
	VisitFunctionStatement(*FunctionStatement)
	VisitReturnStatement(*ReturnStatement)
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

func (s *ExpressionStatement) Accept(v StatementVisitor) {
	v.VisitExpressionStatement(s)
}

type VarStatement struct {
	name        *token.Token
	initializer Expression
}

func NewVarStatement(name *token.Token, initializer Expression) *VarStatement {
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

func (s VarStatement) Name() *token.Token {
	return s.name
}

func (s *VarStatement) Accept(v StatementVisitor) {
	v.VisitVarStatement(s)
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

func (s *BlockStatement) Accept(v StatementVisitor) {
	v.VisitBlockStatement(s)
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

func (s *IfStatement) Accept(v StatementVisitor) {
	v.VisitIfStatement(s)
}

type WhileStatement struct {
	condition Expression
	body      Statement
}

func NewWhileStatement(condition Expression, body Statement) *WhileStatement {
	return &WhileStatement{
		condition: condition,
		body:      body,
	}
}

func (WhileStatement) statement() bool {
	return true
}

func (s WhileStatement) String() string {
	return fmt.Sprintf("while (%s) %s", s.condition.String(), s.body.String())
}

func (s WhileStatement) Condition() Expression {
	return s.condition
}

func (s WhileStatement) Body() Statement {
	return s.body
}

func (s *WhileStatement) Accept(v StatementVisitor) {
	v.VisitWhileStatement(s)
}

type FunctionStatement struct {
	name   *token.Token
	params []*token.Token
	body   []Statement
}

func NewFunctionStatement(name *token.Token, params []*token.Token, body []Statement) *FunctionStatement {
	return &FunctionStatement{
		name:   name,
		params: params,
		body:   body,
	}
}

func (FunctionStatement) statement() bool {
	return true
}

func (s FunctionStatement) String() string {
	paramStrs := []string{}
	for _, param := range s.params {
		paramStrs = append(paramStrs, param.GetLexeme())
	}
	statementStrings := []string{}
	for _, statement := range s.body {
		statementStrings = append(statementStrings, "\t"+statement.String())
	}
	return fmt.Sprintf("fun %s (%s) {\n%s\n}", s.name.GetLexeme(), strings.Join(paramStrs, ", "), strings.Join(statementStrings, "\n"))
}

func (s FunctionStatement) Name() *token.Token {
	return s.name
}

func (s FunctionStatement) Parameters() []*token.Token {
	return s.params
}

func (s FunctionStatement) Body() []Statement {
	return s.body
}

func (s *FunctionStatement) Accept(v StatementVisitor) {
	v.VisitFunctionStatement(s)
}

type ReturnStatement struct {
	keyword *token.Token
	value   Expression
}

func NewReturnStatement(keyword *token.Token, value Expression) *ReturnStatement {
	return &ReturnStatement{
		keyword: keyword,
		value:   value,
	}
}

func (ReturnStatement) statement() bool {
	return true
}

func (s ReturnStatement) String() string {
	return "return " + s.value.String()
}

func (s ReturnStatement) Keyword() *token.Token {
	return s.keyword
}

func (s ReturnStatement) Value() Expression {
	return s.value
}

func (s *ReturnStatement) Accept(v StatementVisitor) {
	v.VisitReturnStatement(s)
}

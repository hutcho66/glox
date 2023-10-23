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
	left, right Expression
	operator    *token.Token
}

func NewBinaryExpression(left Expression, operator *token.Token, right Expression) Expression {
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

func (b BinaryExpression) Operator() *token.Token {
	return b.operator
}

type TernaryExpression struct {
	condition, consequence, alternative Expression
	operator                            *token.Token
}

func NewTernaryExpression(operator *token.Token, condition, consequence, alternative Expression) Expression {
	return &TernaryExpression{
		condition:   condition,
		consequence: consequence,
		alternative: alternative,
		operator:    operator,
	}
}

// TernaryExpression implements Expression
func (TernaryExpression) expression() bool {
	return true
}

func (e TernaryExpression) String() string {
	return fmt.Sprintf("%s ? %s : %s", e.condition, e.consequence, e.alternative)
}

func (e *TernaryExpression) Accept(v ExpressionVisitor) any {
	return v.VisitTernaryExpression(e)
}

func (e TernaryExpression) Operator() *token.Token {
	return e.operator
}

func (e TernaryExpression) Condition() Expression {
	return e.condition
}

func (e TernaryExpression) Consequence() Expression {
	return e.consequence
}

func (e TernaryExpression) Alternative() Expression {
	return e.alternative
}

type LogicalExpression struct {
	left, right Expression
	operator    *token.Token
}

func NewLogicalExpression(left Expression, operator *token.Token, right Expression) Expression {
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

func (b LogicalExpression) Operator() *token.Token {
	return b.operator
}

type UnaryExpression struct {
	expr     Expression
	operator *token.Token
}

func NewUnaryExpression(operator *token.Token, expr Expression) Expression {
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

func (u UnaryExpression) Operator() *token.Token {
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
	name *token.Token
}

func NewVariableExpression(name *token.Token) Expression {
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

func (e VariableExpression) Name() *token.Token {
	return e.name
}

type AssignmentExpression struct {
	name  *token.Token
	value Expression
}

func NewAssignmentExpression(name *token.Token, value Expression) Expression {
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

func (e AssignmentExpression) Name() *token.Token {
	return e.name
}

func (e AssignmentExpression) Value() Expression {
	return e.value
}

type IndexedAssignmentExpression struct {
	left  *IndexExpression
	value Expression
}

func NewIndexedAssignmentExpressionn(left *IndexExpression, value Expression) Expression {
	return &IndexedAssignmentExpression{
		left:  left,
		value: value,
	}
}

// IndexedAssignmentExpression implements Expression
func (IndexedAssignmentExpression) expression() bool {
	return true
}

func (e IndexedAssignmentExpression) String() string {
	return e.left.String() + " = " + e.value.String()
}

func (e *IndexedAssignmentExpression) Accept(v ExpressionVisitor) any {
	return v.VisitIndexedAssignmentExpression(e)
}

func (e IndexedAssignmentExpression) Left() *IndexExpression {
	return e.left
}

func (e IndexedAssignmentExpression) Value() Expression {
	return e.value
}

type CallExpression struct {
	callee       Expression
	arguments    []Expression
	closingParen *token.Token
}

func NewCallExpression(callee Expression, arguments []Expression, closingParen *token.Token) Expression {
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

func (e CallExpression) ClosingParen() *token.Token {
	return e.closingParen
}

func (e CallExpression) Callee() Expression {
	return e.callee
}

func (e CallExpression) Arguments() []Expression {
	return e.arguments
}

type LambdaExpression struct {
	operator *token.Token
	function *FunctionStatement
}

func NewLambdaExpression(operator *token.Token, function *FunctionStatement) Expression {
	return &LambdaExpression{
		operator: operator,
		function: function,
	}
}

// LambdaExpression implements Expression
func (LambdaExpression) expression() bool {
	return true
}

func (e LambdaExpression) String() string {
	paramStrs := []string{}
	for _, param := range e.function.params {
		paramStrs = append(paramStrs, param.GetLexeme())
	}
	statementStrings := []string{}
	for _, statement := range e.function.body {
		statementStrings = append(statementStrings, "\t"+statement.String())
	}
	return fmt.Sprintf("fun (%s) {\n%s\n}", strings.Join(paramStrs, ", "), strings.Join(statementStrings, "\n"))
}

func (e *LambdaExpression) Accept(v ExpressionVisitor) any {
	return v.VisitLambdaExpression(e)
}

func (e LambdaExpression) Operator() *token.Token {
	return e.operator
}

func (e LambdaExpression) Function() *FunctionStatement {
	return e.function
}

type SequenceExpression struct {
	items []Expression
}

func NewSequenceExpression(items []Expression) Expression {
	return &SequenceExpression{
		items: items,
	}
}

// SequenceExpression implements Expression
func (SequenceExpression) expression() bool {
	return true
}

func (e SequenceExpression) String() string {
	items := []string{}
	for _, item := range e.items {
		items = append(items, item.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(items, ", "))
}

func (e *SequenceExpression) Accept(v ExpressionVisitor) any {
	return v.VisitSequenceExpression(e)
}

func (e SequenceExpression) Items() []Expression {
	return e.items
}

type ArrayExpression struct {
	items []Expression
}

func NewArrayExpression(items []Expression) Expression {
	return &ArrayExpression{
		items: items,
	}
}

// ArrayExpression implements Expression
func (ArrayExpression) expression() bool {
	return true
}

func (e ArrayExpression) String() string {
	items := []string{}
	for _, item := range e.items {
		items = append(items, item.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(items, ", "))
}

func (e *ArrayExpression) Accept(v ExpressionVisitor) any {
	return v.VisitArrayExpression(e)
}

func (e ArrayExpression) Items() []Expression {
	return e.items
}

type MapExpression struct {
	openingBrace *token.Token
	keys         []Expression
	values       []Expression
}

func NewMapExpression(keys, values []Expression, openingBrace *token.Token) Expression {
	return &MapExpression{
		openingBrace: openingBrace,
		keys:         keys,
		values:       values,
	}
}

// MapExpression implements Expression
func (MapExpression) expression() bool {
	return true
}

func (e MapExpression) String() string {
	buf := []string{}
	for i := range e.keys {
		buf = append(buf, e.keys[i].String()+": "+e.values[i].String())
	}
	return fmt.Sprintf("{%s}", strings.Join(buf, ", "))
}

func (e *MapExpression) Accept(v ExpressionVisitor) any {
	return v.VisitMapExpression(e)
}

func (e MapExpression) Keys() []Expression {
	return e.keys
}

func (e MapExpression) Values() []Expression {
	return e.values
}

func (e MapExpression) OpeningBrace() *token.Token {
	return e.openingBrace
}

type IndexExpression struct {
	object         Expression
	leftIndex      Expression
	rightIndex     Expression
	closingBracket *token.Token
}

func NewIndexExpression(object, leftIndex, rightIndex Expression, closingBracket *token.Token) Expression {
	return &IndexExpression{
		object:         object,
		leftIndex:      leftIndex,
		rightIndex:     rightIndex,
		closingBracket: closingBracket,
	}
}

// IndexExpression implements Expression
func (IndexExpression) expression() bool {
	return true
}

func (e IndexExpression) String() string {
	if e.rightIndex != nil {
		return fmt.Sprintf("%s[%s:%s]", e.object.String(), e.leftIndex.String(), e.rightIndex.String())
	}
	return fmt.Sprintf("%s[%s]", e.object.String(), e.leftIndex.String())
}

func (e *IndexExpression) Accept(v ExpressionVisitor) any {
	return v.VisitIndexExpression(e)
}

func (e IndexExpression) Object() Expression {
	return e.object
}

func (e IndexExpression) LeftIndex() Expression {
	return e.leftIndex
}

func (e IndexExpression) RightIndex() Expression {
	return e.rightIndex
}

func (e IndexExpression) ClosingBracket() *token.Token {
	return e.closingBracket
}

package resolver

import (
	"github.com/hutcho66/glox/src/pkg/ast"
	"github.com/hutcho66/glox/src/pkg/interpreter"
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/token"
)

type FunctionType string

const (
	NONE     FunctionType = "NONE"
	FUNCTION              = "FUNCTION"
)

type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
	loop            bool
}

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          []map[string]bool{},
		currentFunction: NONE,
		loop:            false,
	}
}

func (r *Resolver) Resolve(statements []ast.Statement) (ok bool) {
	defer func() {
		// catch any errors
		if err := recover(); err != nil {
			ok = false
			return
		}
	}()

	r.resolveStatements(statements)
	return true
}

func (r *Resolver) resolveStatements(statements []ast.Statement) {
	for _, s := range statements {
		r.resolveStatement(s)
	}
}

func (r *Resolver) resolveStatement(statement ast.Statement) {
	statement.Accept(r)
}

func (r *Resolver) resolveExpression(expression ast.Expression) {
	expression.Accept(r)
}

func (r *Resolver) resolveLocal(expression ast.Expression, name *token.Token) {
	for i := range r.scopes {
		i = len(r.scopes) - 1 - i // reverse order!
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.Resolve(expression, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function *ast.FunctionStatement, functionType FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType

	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStatements(function.Body)
	r.endScope()

	r.currentFunction = enclosingFunction
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	// remove last element of scope
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) peekScope() map[string]bool {
	return r.scopes[len(r.scopes)-1]
}

func (r *Resolver) declare(name *token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.peekScope()

	if _, ok := scope[name.Lexeme]; ok {
		panic(lox_error.ResolutionError(name, "Already a variable with this name in scope"))
	}

	scope[name.Lexeme] = false
}

func (r *Resolver) define(name *token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.peekScope()[name.Lexeme] = true
}

// Resolver implements ast.StatementVisitor.
func (r *Resolver) VisitBlockStatement(s *ast.BlockStatement) {
	r.beginScope()
	r.resolveStatements(s.Statements)
	r.endScope()
}

func (r *Resolver) VisitExpressionStatement(e *ast.ExpressionStatement) {
	r.resolveExpression(e.Expr)
}

func (r *Resolver) VisitFunctionStatement(s *ast.FunctionStatement) {
	r.declare(s.Name)
	r.define(s.Name)

	r.resolveFunction(s, FUNCTION)
}

func (r *Resolver) VisitIfStatement(s *ast.IfStatement) {
	r.resolveExpression(s.Condition)
	r.resolveStatement(s.Consequence)
	if s.Alternative != nil {
		r.resolveStatement(s.Alternative)
	}
}

func (r *Resolver) VisitReturnStatement(s *ast.ReturnStatement) {
	if r.currentFunction == NONE {
		lox_error.ResolutionError(s.Keyword, "Can't return from top level code")
	}
	if s.Value != nil {
		r.resolveExpression(s.Value)
	}
}

func (r *Resolver) VisitBreakStatement(s *ast.BreakStatement) {
	if r.loop == false {
		lox_error.ResolutionError(s.Keyword, "Can't break when not in loop")
	}
}

func (r *Resolver) VisitContinueStatement(s *ast.ContinueStatement) {
	if r.loop == false {
		lox_error.ResolutionError(s.Keyword, "Can't continue when not in loop")
	}
}

func (r *Resolver) VisitVarStatement(s *ast.VarStatement) {
	r.declare(s.Name)
	if s.Initializer != nil {
		r.resolveExpression(s.Initializer)
	}
	r.define(s.Name)
}

func (r *Resolver) VisitLoopStatement(s *ast.LoopStatement) {
	r.resolveExpression(s.Condition)

	r.loop = true
	r.resolveStatement(s.Body)
	if s.Increment != nil {
		r.resolveExpression(s.Increment)
	}
	r.loop = false
}

func (r *Resolver) VisitForEachStatement(s *ast.ForEachStatement) {
	r.resolveExpression(s.Array)

	// we need to begin a new scope here to contain the loop variable
	r.beginScope()
	r.declare(s.VariableName)
	r.define(s.VariableName)

	r.loop = true
	r.resolveStatement(s.Body)
	r.loop = false

	r.endScope()
}

// Resolver implements ast.ExprVisitor.
func (r *Resolver) VisitAssignmentExpression(e *ast.AssignmentExpression) any {
	r.resolveExpression(e.Value)
	r.resolveLocal(e, e.Name)
	return nil
}

func (r *Resolver) VisitTernaryExpression(e *ast.TernaryExpression) any {
	r.resolveExpression(e.Condition)
	r.resolveExpression(e.Consequence)
	r.resolveExpression(e.Alternative)
	return nil
}

func (r *Resolver) VisitBinaryExpression(e *ast.BinaryExpression) any {
	r.resolveExpression(e.Left)
	r.resolveExpression(e.Right)
	return nil
}

func (r *Resolver) VisitCallExpression(e *ast.CallExpression) any {
	r.resolveExpression(e.Callee)
	for _, arg := range e.Arguments {
		r.resolveExpression(arg)
	}
	return nil
}

func (r *Resolver) VisitIndexExpression(e *ast.IndexExpression) any {
	r.resolveExpression(e.Object)
	r.resolveExpression(e.LeftIndex)
	if e.RightIndex != nil {
		r.resolveExpression(e.RightIndex)
	}
	return nil
}

func (r *Resolver) VisitArrayExpression(e *ast.ArrayExpression) any {
	for _, item := range e.Items {
		r.resolveExpression(item)
	}
	return nil
}

func (r *Resolver) VisitMapExpression(e *ast.MapExpression) any {
	for i := range e.Keys {
		r.resolveExpression(e.Keys[i])
		r.resolveExpression(e.Values[i])
	}
	return nil
}

func (r *Resolver) VisitIndexedAssignmentExpression(e *ast.IndexedAssignmentExpression) any {
	r.resolveExpression(e.Left)
	r.resolveExpression(e.Value)
	return nil
}

func (r *Resolver) VisitSequenceExpression(e *ast.SequenceExpression) any {
	for _, item := range e.Items {
		r.resolveExpression(item)
	}
	return nil
}

func (r *Resolver) VisitGroupedExpression(e *ast.GroupingExpression) any {
	r.resolveExpression(e.Expr)
	return nil
}

func (r *Resolver) VisitLambdaExpression(e *ast.LambdaExpression) any {
	r.resolveFunction(e.Function, FUNCTION)
	return nil
}

func (r *Resolver) VisitLiteralExpression(e *ast.LiteralExpression) any {
	return nil
}

func (r *Resolver) VisitLogicalExpression(e *ast.LogicalExpression) any {
	r.resolveExpression(e.Left)
	r.resolveExpression(e.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpression(e *ast.UnaryExpression) any {
	r.resolveExpression(e.Expr)
	return nil
}

func (r *Resolver) VisitVariableExpression(e *ast.VariableExpression) any {
	if len(r.scopes) > 0 {
		if val, ok := r.peekScope()[e.Name.Lexeme]; ok && val == false {
			// visiting declared but not yet defined variable is an error
			panic(lox_error.ResolutionError(e.Name, "Can't read local variable in its own initializer"))
		}
	}

	r.resolveLocal(e, e.Name)
	return nil
}

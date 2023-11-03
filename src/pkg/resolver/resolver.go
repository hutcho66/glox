package resolver

import (
	"github.com/hutcho66/glox/src/pkg/ast"
	"github.com/hutcho66/glox/src/pkg/interpreter"
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/token"
)

type FunctionType int
type ClassType int

const (
	NOT_FUNCTION FunctionType = iota
	FUNCTION
	METHOD
	INITIALIZER
)

const (
	NOT_CLASS ClassType = iota
	CLASS
	SUBCLASS
)

type Resolver struct {
	errors          *lox_error.LoxErrors
	interpreter     *interpreter.Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
	currentClass    ClassType
	currentMethod   ast.MethodType
	loop            bool
}

func NewResolver(interpreter *interpreter.Interpreter, errors *lox_error.LoxErrors) *Resolver {
	return &Resolver{
		errors:          errors,
		interpreter:     interpreter,
		scopes:          []map[string]bool{},
		currentFunction: NOT_FUNCTION,
		currentClass:    NOT_CLASS,
		currentMethod:   ast.NOT_METHOD,
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

	if function.Kind == ast.STATIC_METHOD && r.currentClass == NOT_CLASS {
		panic(r.errors.ResolutionError(function.Name, "Cannot declare function as static outside of class declaration."))
	}

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
		panic(r.errors.ResolutionError(name, "Already a variable with this name in scope"))
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

func (r *Resolver) VisitClassStatement(s *ast.ClassStatement) {
	enclosingClass := r.currentClass
	r.currentClass = CLASS

	r.declare(s.Name)
	r.define(s.Name)

	if s.Superclass != nil {
		if s.Name.Lexeme == s.Superclass.Name.Lexeme {
			panic(r.errors.ResolutionError(s.Superclass.Name, "A class can't inherit from itself."))
		}

		r.currentClass = SUBCLASS
		r.resolveExpression(s.Superclass)
	}

	if s.Superclass != nil {
		r.beginScope()
		r.peekScope()["super"] = true
	}

	r.beginScope()

	r.peekScope()["this"] = true
	for _, method := range s.Methods {
		if method.Name.Lexeme == "init" {
			if method.Kind != ast.NORMAL_METHOD {
				panic(r.errors.ResolutionError(s.Name, "init method cannot be static, getter or setter"))
			}
			r.resolveFunction(method, INITIALIZER)
		} else {
			r.currentMethod = method.Kind
			r.resolveFunction(method, METHOD)
			r.currentMethod = ast.NOT_METHOD
		}
	}

	r.endScope()

	if s.Superclass != nil {
		r.endScope()
	}

	r.currentClass = enclosingClass
}

func (r *Resolver) VisitIfStatement(s *ast.IfStatement) {
	r.resolveExpression(s.Condition)
	r.resolveStatement(s.Consequence)
	if s.Alternative != nil {
		r.resolveStatement(s.Alternative)
	}
}

func (r *Resolver) VisitReturnStatement(s *ast.ReturnStatement) {
	if r.currentFunction == NOT_FUNCTION {
		panic(r.errors.ResolutionError(s.Keyword, "Can't return from top level code"))
	}
	if s.Value != nil {
		if r.currentFunction == INITIALIZER {
			panic(r.errors.ResolutionError(s.Keyword, "Can't return a value from an initializer"))
		}
		if r.currentMethod == ast.SETTER_METHOD {
			panic(r.errors.ResolutionError(s.Keyword, "Can't return a value from a setter"))
		}
		r.resolveExpression(s.Value)
	}
}

func (r *Resolver) VisitBreakStatement(s *ast.BreakStatement) {
	if r.loop == false {
		panic(r.errors.ResolutionError(s.Keyword, "Can't break when not in loop"))
	}
}

func (r *Resolver) VisitContinueStatement(s *ast.ContinueStatement) {
	if r.loop == false {
		panic(r.errors.ResolutionError(s.Keyword, "Can't continue when not in loop"))
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

func (r *Resolver) VisitGetExpression(e *ast.GetExpression) any {
	r.resolveExpression(e.Object)
	return nil
}

func (r *Resolver) VisitSetExpression(e *ast.SetExpression) any {
	r.resolveExpression(e.Value)
	r.resolveExpression(e.Object)
	return nil
}

func (r *Resolver) VisitThisExpression(e *ast.ThisExpression) any {
	if r.currentClass == NOT_CLASS {
		panic(r.errors.ResolutionError(e.Keyword, "Can't use 'this' outside of a class."))
	}
	r.resolveLocal(e, e.Keyword)
	return nil
}

func (r *Resolver) VisitSuperGetExpression(e *ast.SuperGetExpression) any {
	if r.currentClass == NOT_CLASS {
		panic(r.errors.ResolutionError(e.Keyword, "Can't use 'super' outside of a class."))
	} else if r.currentClass != SUBCLASS {
		panic(r.errors.ResolutionError(e.Keyword, "Can't use 'super' in a class with no superclass."))
	}
	r.resolveLocal(e, e.Keyword)
	return nil
}

func (r *Resolver) VisitSuperSetExpression(e *ast.SuperSetExpression) any {
	if r.currentClass == NOT_CLASS {
		panic(r.errors.ResolutionError(e.Keyword, "Can't use 'super' outside of a class."))
	} else if r.currentClass != SUBCLASS {
		panic(r.errors.ResolutionError(e.Keyword, "Can't use 'super' in a class with no superclass."))
	}
	r.resolveLocal(e, e.Keyword)
	r.resolveExpression(e.Value)
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
			panic(r.errors.ResolutionError(e.Name, "Can't read local variable in its own initializer"))
		}
	}

	r.resolveLocal(e, e.Name)
	return nil
}

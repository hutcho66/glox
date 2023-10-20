package parser

import (
	"github.com/hutcho66/glox/src/pkg/ast"
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() []ast.Statement {
	statements := []ast.Statement{}
	for !p.isAtEnd() && !p.check(token.NEW_LINE) {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) declaration() (declaration ast.Statement) {
	// catch any panics and synchronize to recover
	defer func() {
		if err := recover(); err != nil {
			p.synchronize()

			// return nil for this statement
			declaration = nil
			return
		}
	}()

	if p.match(token.VAR) {
		declaration = p.varDeclaration()
	} else {
		declaration = p.statement()
	}

	return
}

func (p *Parser) varDeclaration() ast.Statement {
	name := p.consume(token.IDENTIFIER, "Expect variable name.")

	var initializer ast.Expression = nil
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}

	p.endStatement()

	return ast.NewVarStatement(name, initializer)
}

func (p *Parser) statement() ast.Statement {
	if p.match(token.PRINT) {
		return p.printStatement()
	}

	if p.match(token.IF) {
		return p.ifStatement()
	}

	if p.match(token.LEFT_BRACE) {
		return ast.NewBlockStatement(p.block())
	}

	return p.expressionStatement()
}

func (p *Parser) block() []ast.Statement {
	statements := []ast.Statement{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statement := p.declaration()
		statements = append(statements, statement)
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after block")
	return statements
}

func (p *Parser) ifStatement() ast.Statement {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'if'")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after if condition")
	consequence := p.statement()

	var alternative ast.Statement = nil
	if p.match(token.ELSE) {
		alternative = p.statement()
	}

	return ast.NewIfStatement(condition, consequence, alternative)
}

func (p *Parser) printStatement() ast.Statement {
	expr := p.expression()
	p.endStatement()

	return ast.NewPrintStatement(expr)
}

func (p *Parser) expressionStatement() ast.Statement {
	expr := p.expression()
	p.endStatement()
	return ast.NewExpressionStatement(expr)
}

func (p *Parser) expression() ast.Expression {
	return p.assignment()
}

func (p *Parser) assignment() ast.Expression {
	expr := p.or()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if varExpr, ok := expr.(*ast.VariableExpression); ok {
			name := varExpr.Name()
			return ast.NewAssignmentExpression(name, value)
		}

		panic(lox_error.ParserError(equals, "Invalid assignment target"))
	}

	return expr
}

func (p *Parser) or() ast.Expression {
	expr := p.and()

	for p.match(token.OR) {
		operator := p.previous()
		right := p.and()

		expr = ast.NewLogicalExpression(expr, operator, right)
	}

	return expr
}

func (p *Parser) and() ast.Expression {
	expr := p.equality()

	for p.match(token.AND) {
		operator := p.previous()
		right := p.equality()

		expr = ast.NewLogicalExpression(expr, operator, right)
	}

	return expr
}

func (p *Parser) equality() ast.Expression {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()

		expr = ast.NewBinaryExpression(expr, operator, right)
	}

	return expr
}

func (p *Parser) comparison() ast.Expression {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()

		expr = ast.NewBinaryExpression(expr, operator, right)
	}

	return expr
}

func (p *Parser) term() ast.Expression {
	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = ast.NewBinaryExpression(expr, operator, right)
	}

	return expr
}

func (p *Parser) factor() ast.Expression {
	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = ast.NewBinaryExpression(expr, operator, right)
	}

	return expr
}

func (p *Parser) unary() ast.Expression {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.NewUnaryExpression(operator, right)
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expression {
	if p.match(token.FALSE) {
		return ast.NewLiteralExpression(false)
	}
	if p.match(token.TRUE) {
		return ast.NewLiteralExpression(true)
	}
	if p.match(token.NIL) {
		return ast.NewLiteralExpression(nil)
	}
	if p.match(token.NUMBER, token.STRING) {
		return ast.NewLiteralExpression(p.previous().GetLiteral())
	}
	if p.match(token.IDENTIFIER) {
		return ast.NewVariableExpression(p.previous())
	}
	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")

		return ast.NewGroupingExpression(expr)
	}

	panic(lox_error.ParserError(p.peek(), "Expect expression."))
}

func (p *Parser) consume(tokenType token.TokenType, message string) token.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	err := lox_error.ParserError(p.peek(), message)
	panic(err)
}

func (p *Parser) endStatement() {
	// Must have at least one semicolon or newline to terminate a statement
	if terminated := p.match(token.SEMICOLON, token.NEW_LINE); !terminated {
		panic(lox_error.ParserError(p.peek(), "Improperly terminated statement"))
	}

	// Consume as many extra newlines as possible
	for p.match(token.NEW_LINE) {
		continue
	}
}

func (p *Parser) match(tokenTypes ...token.TokenType) bool {
	for _, t := range tokenTypes {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().GetType() == tokenType
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().GetType() == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().GetType() == token.SEMICOLON {
			return
		}

		switch p.peek().GetType() {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}

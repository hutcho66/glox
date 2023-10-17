package parser

import (
	"github.com/hutcho66/glox/src/pkg/ast"
	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/token"
)

type Parser struct {
	tokens []token.Token
	current int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens: tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ast.Expression {
	expr, err := p.expression();
	if err != nil {
		return nil;
	}
	return expr;
}

func (p *Parser) expression() (ast.Expression, error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expression, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err;
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err;
		}
		expr = ast.NewBinaryExpression(expr, operator, right)
	}

	return expr, nil;
}

func (p *Parser) comparison() (ast.Expression, error) {
	expr, err := p.term();
	if err != nil {
		return nil, err;
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous();
		right, err := p.term();
		if err != nil {
			return nil, err;
		}
		expr = ast.NewBinaryExpression(expr, operator, right);
	}

	return expr, nil;
}

func (p *Parser) term() (ast.Expression, error) {
	expr, err := p.factor();
	if err != nil {
		return nil, err;
	}

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous();
		right, err := p.factor();
		if err != nil {
			return nil, err;
		}
		expr = ast.NewBinaryExpression(expr, operator, right);
	}

	return expr, nil;
}

func (p *Parser) factor() (ast.Expression, error) {
	expr, err := p.unary();
	if err != nil {
		return nil, err;
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous();
		right, err := p.unary();
		if err != nil {
			return nil, err;
		}
		expr = ast.NewBinaryExpression(expr, operator, right);
	}

	return expr, nil;
}

func (p *Parser) unary() (ast.Expression, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous();
		right, err := p.unary();
		if err != nil {
			return nil, err;
		}
		return ast.NewUnaryExpression(operator, right), nil;
	}

	return p.primary();
}

func (p *Parser) primary() (ast.Expression, error) {
	if p.match(token.FALSE) {
		return ast.NewLiteralExpression(false), nil;
	}
	if p.match(token.TRUE) {
		return ast.NewLiteralExpression(true), nil;
	}
	if p.match(token.NIL) {
		return ast.NewLiteralExpression(nil), nil;
	}
	if p.match(token.NUMBER, token.STRING) {
		return ast.NewLiteralExpression(p.previous().GetLiteral()), nil;
	}
	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression();
		if err != nil {
			return nil, err;
		}
		if _, err := p.consume(token.RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return nil, err
		}
		return ast.NewGroupingExpression(expr), nil;
	}

	return nil, lox_error.ParserError(p.peek(), "Expect expression.")
}

func (p *Parser) consume(tokenType token.TokenType, message string) (token.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil;
	}

	return token.Token{}, lox_error.ParserError(p.peek(), message);
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
	return p.tokens[p.current - 1]
}

// func (p *Parser) synchronize() {
// 	p.advance();

// 	for !p.isAtEnd() {
// 		if p.previous().GetType() == token.SEMICOLON {
// 			return;
// 		}

// 		switch (p.peek().GetType()) {
// 			case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN: return;
// 		}

// 		p.advance();
// 	}
// }

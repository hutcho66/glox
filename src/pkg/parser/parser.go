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
	for !p.isAtEnd() {
		if !p.match(token.NEW_LINE) {
			statements = append(statements, p.declaration())
		}
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
		return p.varDeclaration()
	} else if p.match(token.FUN) {
		return p.funDeclaration("function")
	} else {
		return p.statement()
	}

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

func (p *Parser) funDeclaration(kind string) ast.Statement {
	name := p.consume(token.IDENTIFIER, "Expect "+kind+" name")
	p.consume(token.LEFT_PAREN, "Expect '(' after "+kind+" name")
	parameters := []*token.Token{}
	if !p.check(token.RIGHT_PAREN) {
		for ok := true; ok; ok = p.match(token.COMMA) {
			if len(parameters) >= 255 {
				panic(lox_error.ParserError(p.peek(), "Can't have more than 255 parameters"))
			}

			parameters = append(parameters, p.consume(token.IDENTIFIER, "Expect parameter name"))
		}
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after parameters")

	p.consume(token.LEFT_BRACE, "Expect '{' before "+kind+" body")
	body := p.block()

	return ast.NewFunctionStatement(name, parameters, body)
}

func (p *Parser) statement() ast.Statement {
	if p.match(token.NEW_LINE) {
		// consume and retry
		return p.statement()
	}

	if p.match(token.RETURN) {
		return p.returnStatement()
	}

	if p.match(token.BREAK) {
		return p.breakStatement()
	}

	if p.match(token.CONTINUE) {
		return p.continueStatement()
	}

	if p.match(token.IF) {
		return p.ifStatement()
	}

	if p.match(token.WHILE) {
		return p.whileStatement()
	}

	if p.match(token.FOR) {
		return p.forStatement()
	}

	if p.check(token.LEFT_BRACE) {
		if p.checkAhead(token.RIGHT_BRACE, 1) || (p.checkAhead(token.STRING, 1) && !p.checkAhead(token.COLON, 2)) {
			// this looks like a map
			return p.expressionStatement()
		}
		p.match(token.LEFT_BRACE)
		return ast.NewBlockStatement(p.block())
	}

	return p.expressionStatement()
}

func (p *Parser) block() []ast.Statement {
	statements := []ast.Statement{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		if !p.match(token.NEW_LINE) {
			statement := p.declaration()
			statements = append(statements, statement)
		}
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after block")
	return statements
}

func (p *Parser) returnStatement() ast.Statement {
	keyword := p.previous()
	var value ast.Expression = nil
	if !p.check(token.SEMICOLON) && !p.check(token.NEW_LINE) {
		value = p.expression()
	}

	p.endStatement()
	return ast.NewReturnStatement(keyword, value)
}

func (p *Parser) breakStatement() ast.Statement {
	keyword := p.previous()
	p.endStatement()
	return ast.NewBreakStatement(keyword)
}

func (p *Parser) continueStatement() ast.Statement {
	keyword := p.previous()
	p.endStatement()
	return ast.NewContinueStatement(keyword)
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

func (p *Parser) whileStatement() ast.Statement {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'while'")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after while condition")

	body := p.statement()

	// while statements have no increment
	return ast.NewLoopStatement(condition, body, nil)
}

func (p *Parser) forStatement() ast.Statement {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'for'")

	if p.check(token.VAR) && p.checkAhead(token.OF, 2) {
		// for (IDENT of ARRAY)format
		p.match(token.VAR)
		name := p.consume(token.IDENTIFIER, "Expect variable name after var")
		p.consume(token.OF, "Expect 'of' after variable name")
		array := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses")

		body := p.statement()

		return ast.NewForEachStatement(name, array, body)
	}

	// else continue with c-style loop
	var initializer ast.Statement
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition ast.Expression = nil
	if !p.check(token.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(token.SEMICOLON, "Expect ';' after loop condition")

	var increment ast.Expression = nil
	if !p.check(token.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses")

	body := p.statement()

	if condition == nil {
		// if there is no condition, set it to 'true' to make infinite loop
		condition = ast.NewLiteralExpression(true)
	}
	// create LoopStatement using condition, body and increment
	body = ast.NewLoopStatement(condition, body, increment)

	// if there is an initializer, add before loop statement
	if initializer != nil {
		body = ast.NewBlockStatement([]ast.Statement{
			initializer,
			body,
		})
	}

	return body
}

func (p *Parser) expressionStatement() ast.Statement {
	expr := p.expression()
	p.endStatement()
	return ast.NewExpressionStatement(expr)
}

func (p *Parser) expression() ast.Expression {
	if p.check(token.LEFT_PAREN) {
		// need to check ahead to test if this is a lambda
		if p.checkAhead(token.RIGHT_PAREN, 1) {
			// must be lambda with no params
			return p.lambda()
		}

		if p.checkAhead(token.IDENTIFIER, 1) {
			// presence of comma indicates a lambda
			// as does a right paren and then the arrow operator
			if p.checkAhead(token.COMMA, 2) || p.checkAhead(token.RIGHT_PAREN, 2) && p.checkAhead(token.LAMBDA_ARROW, 3) {
				return p.lambda()
			}
		}
	}

	if p.check(token.IDENTIFIER) && p.checkAhead(token.LAMBDA_ARROW, 1) {
		// x => <expression>
		return p.lambda()
	}

	return p.ternary()
}

func (p *Parser) lambda() ast.Expression {
	parameters := []*token.Token{}
	if p.match(token.IDENTIFIER) {
		// x => <expression> form
		parameters = append(parameters, p.previous())
	} else {
		p.consume(token.LEFT_PAREN, "unexpected error") // already checked

		if !p.check(token.RIGHT_PAREN) {
			for ok := true; ok; ok = p.match(token.COMMA) {
				if len(parameters) >= 255 {
					panic(lox_error.ParserError(p.peek(), "Can't have more than 255 parameters"))
				}

				parameters = append(parameters, p.consume(token.IDENTIFIER, "Expect parameter name"))
			}
		}

		p.consume(token.RIGHT_PAREN, "Expect ')' after parameters")
	}

	operator := p.consume(token.LAMBDA_ARROW, "Expect '=>' after lambda parameters")

	var body []ast.Statement
	if !p.check(token.LEFT_BRACE) || (p.checkAhead(token.STRING, 1) && p.checkAhead(token.COLON, 2)) {
		// this is an expression return lambda
		line := p.peek().GetLine()
		expression := p.expression()
		// add implicit return statement
		token := token.NewToken(token.RETURN, "return", nil, line)
		body = []ast.Statement{
			ast.NewReturnStatement(token, expression),
		}
	} else {
		// this is a block lambda
		p.match(token.LEFT_BRACE)
		body = p.block()
	}

	function := ast.NewFunctionStatement(nil, parameters, body)

	return ast.NewLambdaExpression(operator, function)
}

func (p *Parser) ternary() ast.Expression {
	condition := p.assignment()

	if p.match(token.QUESTION) {
		operator := p.previous()
		consequence := p.expression()
		p.consume(token.COLON, "Expect ':' after expression following '?'")
		alternative := p.expression()

		return ast.NewTernaryExpression(operator, condition, consequence, alternative)
	}

	return condition
}

func (p *Parser) assignment() ast.Expression {
	expr := p.or()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		switch e := expr.(type) {
		case *ast.VariableExpression:
			{
				return ast.NewAssignmentExpression(e.Name(), value)
			}
		case *ast.IndexExpression:
			{
				if e.RightIndex() != nil {
					panic(lox_error.ParserError(equals, "Cannot assign to array slice"))
				}
				return ast.NewIndexedAssignmentExpressionn(e, value)
			}
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

	return p.call_index()
}

func (p *Parser) call_index() ast.Expression {
	expr := p.primary()

	for {
		if p.match(token.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(token.LEFT_BRACKET) {
			expr = p.finishIndex(expr)
		} else {
			break
		}
	}

	return expr
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
		if p.match(token.RIGHT_PAREN) {
			// empty sequence expression
			return ast.NewSequenceExpression([]ast.Expression{})
		}
		exprs := p.expressionList()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression")

		if len(exprs) == 1 {
			return ast.NewGroupingExpression(exprs[0])
		} else {
			return ast.NewSequenceExpression(exprs)
		}
	}
	if p.match(token.LEFT_BRACKET) {
		if p.match(token.RIGHT_BRACKET) {
			// empty array
			return ast.NewArrayExpression([]ast.Expression{})
		}
		exprs := p.expressionList()
		p.consume(token.RIGHT_BRACKET, "Expect ']' after array literal")

		return ast.NewArrayExpression(exprs)
	}
	if p.match(token.LEFT_BRACE) {
		openingBrace := p.previous()
		// eat any newlines, they are allowed before first key-pair
		p.eatNewLines()

		if p.match(token.RIGHT_BRACE) {
			// empty array
			return ast.NewMapExpression([]ast.Expression{}, []ast.Expression{}, openingBrace)
		}

		keys := []ast.Expression{}
		values := []ast.Expression{}
		for ok := true; ok; ok = p.match(token.COMMA) {
			p.eatNewLines()

			keys = append(keys, p.expression())
			p.consume(token.COLON, "Expect ':' between key and value in map literal")
			values = append(values, p.expression())

			p.eatNewLines()
		}
		p.consume(token.RIGHT_BRACE, "Expect '}' after map literal")

		return ast.NewMapExpression(keys, values, openingBrace)
	}

	panic(lox_error.ParserError(p.peek(), "Expect expression."))
}

func (p *Parser) expressionList() []ast.Expression {
	// eat any newlines, they are allowed before first expression in list
	p.eatNewLines()

	exprs := []ast.Expression{}
	for ok := true; ok; ok = p.match(token.COMMA) {
		p.eatNewLines()
		exprs = append(exprs, p.expression())
		p.eatNewLines()
	}
	return exprs
}

func (p *Parser) finishIndex(array ast.Expression) ast.Expression {
	leftIndex := p.expression()
	var rightIndex ast.Expression
	if p.match(token.COLON) {
		rightIndex = p.expression()
	}
	closingBracket := p.consume(token.RIGHT_BRACKET, "Expect ']' after index")

	return ast.NewIndexExpression(array, leftIndex, rightIndex, closingBracket)
}

func (p *Parser) finishCall(callee ast.Expression) ast.Expression {
	args := []ast.Expression{}
	if !p.check(token.RIGHT_PAREN) {
		for ok := true; ok; ok = p.match(token.COMMA) {
			if len(args) >= 255 {
				panic(lox_error.ParserError(p.peek(), "Can't have more than 255 arguments"))
			}
			args = append(args, p.expression())
		}
	}
	closingParen := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments")

	return ast.NewCallExpression(callee, args, closingParen)
}

func (p *Parser) consume(tokenType token.TokenType, message string) *token.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	err := lox_error.ParserError(p.peek(), message)
	panic(err)
}

func (p *Parser) endStatement() {
	// a closing brace is a valid statement ending, to allow statements like this on one line
	// `if (true) { var x = 5; print(x) }`
	if p.check(token.RIGHT_BRACE) {
		return
	}

	// Otherwise, must have at least one semicolon or newline to terminate a statement
	if terminated := p.match(token.SEMICOLON, token.NEW_LINE); !terminated && !p.isAtEnd() {
		panic(lox_error.ParserError(p.peek(), "Improperly terminated statement"))
	}

	// Consume as many extra newlines as possible
	p.eatNewLines()
}

func (p *Parser) eatNewLines() {
	for p.match(token.NEW_LINE) {
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

func (p *Parser) checkAhead(tokenType token.TokenType, lookahead int) bool {
	position := p.current + lookahead
	if position >= len(p.tokens) {
		return false
	}
	return p.tokens[position].GetType() == tokenType
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().GetType() == token.EOF
}

func (p *Parser) peek() *token.Token {
	return &p.tokens[p.current]
}

func (p *Parser) previous() *token.Token {
	return &p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().GetType() == token.SEMICOLON {
			return
		}

		switch p.peek().GetType() {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.RETURN:
			return
		}

		p.advance()
	}
}

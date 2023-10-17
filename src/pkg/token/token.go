package token

type Token struct {
	tokenType TokenType
	lexeme string
	literal any
	line int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) *Token {
	return &Token{
		tokenType: tokenType,
		lexeme: lexeme,
		literal: literal,
		line: line,
	};
}

func (t Token) GetType() TokenType {
	return t.tokenType;
}

func (t Token) GetLiteral() any {
	return t.literal;
}

func (t Token) GetLexeme() string {
	return t.lexeme;
}

func (t Token) GetLine() int {
	return t.line;
}


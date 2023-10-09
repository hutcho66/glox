package token

import "fmt"

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

func (t Token) String() string {
	if t.literal == nil {
		return fmt.Sprintf("%s %s", t.tokenType, t.lexeme);
	}
	return fmt.Sprintf("%s %s %s", t.tokenType, t.lexeme, t.literal);
}

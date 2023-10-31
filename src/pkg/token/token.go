package token

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

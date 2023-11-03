package token

var keywords = map[string]TokenType{
	"and":      AND,
	"break":    BREAK,
	"class":    CLASS,
	"continue": CONTINUE,
	"else":     ELSE,
	"false":    FALSE,
	"for":      FOR,
	"fun":      FUN,
	"if":       IF,
	"nil":      NIL,
	"of":       OF,
	"or":       OR,
	"return":   RETURN,
	"static":   STATIC,
	"super":    SUPER,
	"this":     THIS,
	"true":     TRUE,
	"var":      VAR,
	"while":    WHILE,
}

func LookupKeyword(word string) TokenType {
	if tokenType, ok := keywords[word]; ok {
		return tokenType
	}
	return IDENTIFIER
}

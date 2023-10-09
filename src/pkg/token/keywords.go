package token

var keywords = map[string]TokenType{
	"and": 		AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
};

func LookupKeyword(word string) TokenType {
	if tokenType, ok := keywords[word]; ok {
		return tokenType;
	}
	return IDENTIFIER;
}

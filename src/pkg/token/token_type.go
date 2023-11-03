package token

type TokenType string

const (
	EOF TokenType = "EOF"

	NEW_LINE = "NEWLINE"

	// Single character symbols
	LEFT_PAREN    = "("
	RIGHT_PAREN   = ")"
	LEFT_BRACE    = "{"
	RIGHT_BRACE   = "}"
	LEFT_BRACKET  = "["
	RIGHT_BRACKET = "]"
	COMMA         = ","
	DOT           = "."
	MINUS         = "-"
	PLUS          = "+"
	SEMICOLON     = ";"
	SLASH         = "/"
	STAR          = "*"
	QUESTION      = "?"
	COLON         = ":"

	// Multi character symbols
	BANG          = "!"
	BANG_EQUAL    = "!="
	EQUAL         = "="
	EQUAL_EQUAL   = "=="
	GREATER       = ">"
	GREATER_EQUAL = ">="
	LESS          = "<"
	LESS_EQUAL    = "<="
	LAMBDA_ARROW  = "=>"

	// Literals
	IDENTIFIER = "IDENTIFIER"
	STRING     = "STRING"
	NUMBER     = "NUMBER"

	// Keywords
	AND      = "AND"
	BREAK    = "BREAK"
	CLASS    = "CLASS"
	CONTINUE = "CONTINUE"
	ELSE     = "ELSE"
	FALSE    = "FALSE"
	FUN      = "FUN"
	FOR      = "FOR"
	GET      = "GET"
	IF       = "IF"
	NIL      = "NIL"
	OF       = "OF"
	OR       = "OR"
	RETURN   = "RETURN"
	SET      = "SET"
	STATIC   = "STATIC"
	SUPER    = "SUPER"
	THIS     = "THIS"
	TRUE     = "TRUE"
	VAR      = "VAR"
	WHILE    = "WHILE"
)

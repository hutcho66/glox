package token

type TokenType string;

const (
  EOF TokenType = "EOF"

  // Single character symbols
  LEFT_PAREN  = "("
  RIGHT_PAREN = ")"
  LEFT_BRACE  = "{"
  RIGHT_BRACE = "}"
  COMMA       = ","
  DOT         = "."
  MINUS       = "-"
  PLUS        = "+"
  SEMICOLON   = ";"
  SLASH       = "/"
  STAR        = "*"

  // Multi character symbols
  BANG          = "!"
  BANG_EQUAL    = "!="
  EQUAL         = "="
  EQUAL_EQUAL   = "=="
  GREATER       = ">"
  GREATER_EQUAL = ">="
  LESS          = "<"
  LESS_EQUAL    = "<="

  // Literals
  IDENTIFIER  = "IDENTIFIER"
  STRING      = "STRING"
  NUMBER      = "NUMBER"

  // Keywords
  AND     = "AND"
  CLASS   = "CLASS"
  ELSE    = "ELSE"
  FALSE   = "FALSE"
  FUN     = "FUN"
  FOR     = "FOR"
  IF      = "IF"
  NIL     = "NIL"
  OR      = "OR"
  PRINT   = "PRINT"
  RETURN  = "RETURN"
  SUPER   = "SUPER"
  THIS    = "THIS"
  TRUE    = "TRUE"
  VAR     = "VAR"
  WHILE   = "WHILE"
)

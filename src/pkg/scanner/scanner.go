package scanner

import (
	"strconv"

	"github.com/hutcho66/glox/src/pkg/lox_error"
	"github.com/hutcho66/glox/src/pkg/token"
)

type Scanner struct {
	source               string
	tokens               []token.Token
	start, current, line int
}

// Public methods

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []token.Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, *token.NewToken(token.EOF, "", nil, s.line))
	return s.tokens
}

// Private methods

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	// Symbols
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)
	case '!':
		s.addTokenConditional('=', token.BANG_EQUAL, token.BANG)
	case '=':
		{
			if s.match('=') {
				s.addToken(token.EQUAL_EQUAL)
			} else if s.match('>') {
				s.addToken(token.LAMBDA_ARROW)
			} else {
				s.addToken(token.EQUAL)
			}
		}
	case '<':
		s.addTokenConditional('=', token.LESS_EQUAL, token.LESS)
	case '>':
		s.addTokenConditional('=', token.GREATER_EQUAL, token.GREATER)
	case '/':
		{
			if s.match('/') {
				// Comment goes to the end of the line
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				s.addToken(token.SLASH)
			}
		}

	// Ignore whitespace
	case ' ':
		break
	case '\r':
		break
	case '\t':
		break

	// Newlines are a token
	case '\n':
		{
			s.addToken(token.NEW_LINE)
			s.line++
		}

	// Literals
	case '"':
		s.string()

	default:
		{
			if isDigit(c) {
				s.number()
			} else if isAlpha(c) {
				s.identifier()
			} else {
				lox_error.ScannerError(s.line, "Unexpected character.")
			}
		}
	}
}

func (s *Scanner) string() {
	// Advance until either EOF or closing quote, incrementing line count when necessary
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		lox_error.ScannerError(s.line, "Unterminated string.")
		return
	}

	// consume closing quote
	s.advance()

	// trim quote symbols
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(token.STRING, value)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addTokenWithLiteral(token.NUMBER, value)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	word := s.source[s.start:s.current]
	s.addToken(token.LookupKeyword(word))
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
}

func (s *Scanner) advance() byte {
	ch := s.source[s.current]
	s.current++
	return ch
}

func (s *Scanner) addToken(tokenType token.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenConditional(expected byte, matchType, elseType token.TokenType) {
	if s.match(expected) {
		s.addToken(matchType)
	} else {
		s.addToken(elseType)
	}
}

func (s *Scanner) addTokenWithLiteral(tokenType token.TokenType, literal any) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, *token.NewToken(tokenType, lexeme, literal, s.line))
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_'
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

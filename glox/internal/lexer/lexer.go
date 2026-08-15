package lexer

import (
	"fmt"
	"strconv"

	errrp "github.com/koftamainee/lox/glox/internal/error"
	"github.com/koftamainee/lox/glox/internal/token"
)

type Lexer struct {
	source string
	tokens []token.Token
	errors errrp.ErrorReporter

	start   int
	current int
	line    int
}

func New(source string, errors errrp.ErrorReporter) Lexer {
	return Lexer{
		source:  source,
		tokens:  make([]token.Token, 0),
		errors:  errors,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (l *Lexer) ScanTokens() []token.Token {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}

	l.tokens = append(l.tokens, token.New(token.EOF, "", nil, l.line))

	return l.tokens
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) scanToken() {
	c := l.advance()
	switch c {
	case '(':
		l.addToken(token.LeftParen)
	case ')':
		l.addToken(token.RightParen)
	case '{':
		l.addToken(token.LeftBrace)
	case '}':
		l.addToken(token.RightBrace)
	case ',':
		l.addToken(token.Comma)
	case '.':
		l.addToken(token.Dot)
	case '-':
		l.addToken(token.Minus)
	case '+':
		l.addToken(token.Plus)
	case ';':
		l.addToken(token.Semicolon)
	case '*':
		l.addToken(token.Star)
	case '!':
		if l.match('=') {
			l.addToken(token.BangEqual)
		} else {
			l.addToken(token.Bang)
		}
	case '=':
		if l.match('=') {
			l.addToken(token.EqualEqual)
		} else {
			l.addToken(token.Equal)
		}
	case '>':
		if l.match('=') {
			l.addToken(token.GreaterEqual)
		} else {
			l.addToken(token.Greater)
		}
	case '<':
		if l.match('=') {
			l.addToken(token.LessEqual)
		} else {
			l.addToken(token.Less)
		}
	case '/':
		if l.match('/') {
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
		} else if l.match('*') {
			l.multilineComment()
		} else {
			l.addToken(token.Slash)
		}

	case ' ':
		break
	case '\r':
		break
	case '\t':
		break
	case '\n':
		l.line++

	case '"':
		l.string()

	default:
		if isDigit(c) {
			l.number()
		} else if isAlpha(c) {
			l.identifier()
		} else {
			l.errors.Error(l.line, "Unexpected character.")
		}
	}
}

func (l *Lexer) string() {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}

	if l.isAtEnd() {
		l.errors.Error(l.line, "Unterminated string")
		return
	}

	// the closing "
	l.advance()

	str_value := l.source[l.start+1 : l.current-1]
	l.addTokenWithLiteral(token.String, str_value)
}

func (l *Lexer) number() {
	for isDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance()
	}

	for isDigit(l.peek()) {
		l.advance()
	}

	strNum := l.source[l.start:l.current]
	doubleNum, err := strconv.ParseFloat(strNum, 64)
	if err != nil {
		panic(fmt.Sprintf("unreachable: invalid number %q", strNum))
	}

	l.addTokenWithLiteral(token.Number, doubleNum)
}

func (l *Lexer) identifier() {
	for isAlphaNumeric(l.peek()) {
		l.advance()
	}

	text := l.source[l.start:l.current]
	tokenType, ok := keywords[text]
	if !ok {
		l.addToken(token.Identifier)
	} else {
		l.addToken(tokenType)
	}
}

func (l *Lexer) multilineComment() {
	depth := 1

	for depth > 0 {
		if l.isAtEnd() {
			l.errors.Error(l.line, "Unterminated comment.")
			return
		}
		if l.peek() == '/' && l.peekNext() == '*' {
			l.advance()
			l.advance()
			depth++
		} else if l.peek() == '*' && l.peekNext() == '/' {
			l.advance()
			l.advance()
			depth--
		} else {
			if l.advance() == '\n' {
				l.line++
			}
		}
	}
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return l.source[l.current+1]
}

func (l *Lexer) match(expected byte) bool {
	if l.isAtEnd() {
		return false
	}

	if l.source[l.current] != expected {
		return false
	}

	l.current++
	return true
}

func (l *Lexer) advance() byte {
	c := l.source[l.current]
	l.current++
	return c
}

func (l *Lexer) addToken(tokenType token.Type) {
	l.addTokenWithLiteral(tokenType, nil)
}

func (l *Lexer) addTokenWithLiteral(tokenType token.Type, literal any) {
	text := l.source[l.start:l.current]
	l.tokens = append(l.tokens, token.New(tokenType, text, literal, l.line))
}

package token

import "fmt"

type Token struct {
	TokenType Type
	Lexeme    string
	Literal   any
	Line      int
}

func New(token_type Type, lexeme string, literal any, line int) Token {
	return Token{
		TokenType: token_type,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%v %s %v", t.TokenType, t.Lexeme, t.Literal)
}

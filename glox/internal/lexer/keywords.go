package lexer

import "github.com/koftamainee/lox/glox/internal/token"

var keywords = map[string]token.Type{
	"and":    token.And,
	"class":  token.Class,
	"else":   token.Else,
	"false":  token.False,
	"fun":    token.Fun,
	"for":    token.For,
	"if":     token.If,
	"nil":    token.Nil,
	"or":     token.Or,
	"print":  token.Print,
	"return": token.Return,
	"super":  token.Super,
	"this":   token.This,
	"true":   token.True,
	"var":    token.Var,
	"while":  token.While,
}

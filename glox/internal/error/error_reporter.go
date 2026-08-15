package error

import "github.com/koftamainee/lox/glox/internal/token"

type ErrorReporter interface {
	Error(line int, message string)
	ErrorAt(token token.Token, message string)
	InternalError(message string)
}

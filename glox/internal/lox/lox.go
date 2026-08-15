package lox

import (
	"fmt"
	"os"

	"github.com/koftamainee/lox/glox/internal/lexer"
	"github.com/koftamainee/lox/glox/internal/parser"
	"github.com/koftamainee/lox/glox/internal/token"
)

type LoxErrorReporter struct {
	HadError bool
}

type Lox struct {
	ErrorReporter LoxErrorReporter
}

func New() (Lox, error) {
	return Lox{
		ErrorReporter: LoxErrorReporter{
			HadError: false}}, nil
}

func (l *Lox) Run(bytes string) error {
	lex := lexer.New(bytes, &l.ErrorReporter)
	tokens := lex.ScanTokens()

	parse := parser.New(tokens, &l.ErrorReporter)
	expr := parse.Parse()

	fmt.Println(expr)

	return nil
}

func (r *LoxErrorReporter) Error(line int, message string) {
	r.report(line, "", message)
}

func (r *LoxErrorReporter) ErrorAt(t token.Token, message string) {
	if t.TokenType == token.EOF {
		r.report(t.Line, "at end", message)
	} else {
		r.report(t.Line, fmt.Sprintf("at '%s'", t.Lexeme), message)
	}
}

func (r *LoxErrorReporter) InternalError(message string) {
	fmt.Fprintf(os.Stderr, "Glox internal error: %s", message)

	r.HadError = true
}

func (r *LoxErrorReporter) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)

	r.HadError = true
}

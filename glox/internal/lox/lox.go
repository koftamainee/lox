package lox

import (
	"fmt"
	"os"

	"github.com/koftamainee/lox/glox/internal/lexer"
	"github.com/koftamainee/lox/glox/internal/parser"
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

	parse := parser.New(tokens)
	expr := parse.Parse()

	fmt.Println(expr)

	return nil
}

func (r *LoxErrorReporter) Error(line int, message string) {
	r.report(line, "", message)
}

func (r *LoxErrorReporter) InternalError(message string) {
	fmt.Fprintf(os.Stderr, "Glox internal error: %s", message)

	r.HadError = true
}

func (r *LoxErrorReporter) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)

	r.HadError = true
}

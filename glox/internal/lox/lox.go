package lox

import (
	"fmt"
	"os"

	"github.com/koftamainee/lox/glox/internal/interpreter"
	"github.com/koftamainee/lox/glox/internal/lexer"
	"github.com/koftamainee/lox/glox/internal/parser"
	"github.com/koftamainee/lox/glox/internal/token"
)

type LoxErrorReporter struct {
	HadError        bool
	HadRuntimeError bool
}

type Lox struct {
	ErrorReporter *LoxErrorReporter

	interpreter interpreter.Interpreter
}

func New() Lox {
	err := LoxErrorReporter{
		HadError:        false,
		HadRuntimeError: false,
	}
	interpreter := interpreter.New(&err)
	return Lox{
		ErrorReporter: &err,
		interpreter:   interpreter,
	}

}

func (l *Lox) Run(bytes string) {
	lex := lexer.New(bytes, l.ErrorReporter)
	tokens := lex.ScanTokens()

	parse := parser.New(tokens, l.ErrorReporter)
	statements := parse.Parse()

	if l.ErrorReporter.HadError {
		return
	}

	l.interpreter.Interpret(statements)
}

func (l *Lox) EvaluateExpression(bytes string) {
	lex := lexer.New(bytes, l.ErrorReporter)
	tokens := lex.ScanTokens()

	parse := parser.New(tokens, l.ErrorReporter)
	expr := parse.ParseExpression()

	if l.ErrorReporter.HadError {
		return
	}

	value := l.interpreter.Evaluate(expr)

	if l.ErrorReporter.HadRuntimeError {
		return
	}

	fmt.Println(value)

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
	fmt.Fprintf(os.Stderr, "Glox internal error: %s\n", message)

	r.HadError = true
}

func (r *LoxErrorReporter) RuntimeError(t token.Token, message string) {
	if t.TokenType == token.EOF {
		r.report(t.Line, "at end", message)
	} else {
		r.report(t.Line, fmt.Sprintf("at '%s'", t.Lexeme), message)
	}
	r.HadRuntimeError = true
}

func (r *LoxErrorReporter) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)

	r.HadError = true
}

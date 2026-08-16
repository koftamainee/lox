package interpreter

import (
	"fmt"

	"github.com/koftamainee/lox/glox/internal/ast"
	"github.com/koftamainee/lox/glox/internal/environment"
)

// type LoxCallable interface {
// 	Call(interpreter *Interpreter, args []any) any
// 	Arity() int
// 	String() string
// }

type loxFunction struct {
	Declaration *ast.FunStatement
}

func (f loxFunction) Call(interpreter *Interpreter, args []any) any {
	env := environment.New(interpreter.globals)

	for i, param := range f.Declaration.Params {
		env.Define(param.Lexeme, args[i])
	}

	block := ast.BlockStatement{Statements: f.Declaration.Body}

	err := interpreter.execBlockStmt(&block, env)
	if err != nil {
		return nil // TODO: proper error handling and value return
	}

	return nil
}

func (f loxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f loxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}

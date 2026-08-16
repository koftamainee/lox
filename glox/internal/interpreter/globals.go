package interpreter

import (
	"bufio"
	"os"
	"time"

	"github.com/koftamainee/lox/glox/internal/environment"
)

func defineGlobals(env *environment.Environment) {
	env.Define("clock", clockFn())
	env.Define("input", inputFn())
}

type nativeFn struct {
	call  func(*Interpreter, []any) any
	arity int
}

func (f nativeFn) Call(interpreter *Interpreter, args []any) any {
	return f.call(interpreter, args)
}

func (f nativeFn) Arity() int {
	return f.arity
}

func (f nativeFn) String() string {
	return "<native fn>"
}

func clockFn() loxCallable {
	call := func(interpreter *Interpreter, args []any) any {
		return float64(time.Now().UnixMilli()) / 1000.0
	}
	return &nativeFn{
		call:  call,
		arity: 0,
	}
}

func inputFn() loxCallable {
	call := func(interpreter *Interpreter, args []any) any {
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return nil
		}

		return scanner.Text()
	}
	return &nativeFn{
		call:  call,
		arity: 0,
	}
}

package interpreter

import "fmt"

type returnErr struct {
	Value any
}

func (e returnErr) Error() string {
	return fmt.Sprint(e.Value)
}

type loxCallable interface {
	Call(interpreter *Interpreter, args []any) (any, error)
	Arity() int
	String() string
}

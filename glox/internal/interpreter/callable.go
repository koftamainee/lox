package interpreter

type loxCallable interface {
	Call(interpreter *Interpreter, args []any) any
	Arity() int
	String() string
}

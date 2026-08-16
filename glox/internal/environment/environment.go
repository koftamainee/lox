package environment

import (
	"errors"

	"github.com/koftamainee/lox/glox/internal/token"
)

var ErrUndeclared = errors.New("variable is not declared")
var ErrUndefined = errors.New("variable is declared, but not defined")

type envBucket struct {
	Value         any
	IsInitialized bool
}

type Environment struct {
	values map[string]envBucket
}

func New() Environment {
	return Environment{
		values: make(map[string]envBucket),
	}
}

func (e *Environment) Declare(name string) {
	e.values[name] = envBucket{Value: nil, IsInitialized: false}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = envBucket{Value: value, IsInitialized: true}
}

func (e *Environment) Assign(name token.Token, value any) error {
	v, ok := e.values[name.Lexeme]
	if !ok {
		return ErrUndeclared
	}

	v.IsInitialized = true
	v.Value = value

	e.values[name.Lexeme] = v
	return nil
}

func (e *Environment) Get(name token.Token) (any, error) {
	v, ok := e.values[name.Lexeme]
	if !ok {
		return nil, ErrUndeclared
	}
	if !v.IsInitialized {
		return nil, ErrUndefined
	}

	return v.Value, nil
}

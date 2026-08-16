package ast

import "github.com/koftamainee/lox/glox/internal/token"

type Statement interface {
	statementImpl()
}

type ExpressionStatement struct {
	Expr Expression
}

type PrintStatement struct {
	Expr Expression
}

type VarStatement struct {
	Name        token.Token
	Initializer Expression
}

func (e *ExpressionStatement) statementImpl() {}
func (e *PrintStatement) statementImpl()      {}
func (e *VarStatement) statementImpl()        {}

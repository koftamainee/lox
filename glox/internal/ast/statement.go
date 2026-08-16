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

type BlockStatement struct {
	Statements []Statement
}

func (e *ExpressionStatement) statementImpl() {}
func (e *PrintStatement) statementImpl()      {}
func (e *VarStatement) statementImpl()        {}
func (e *BlockStatement) statementImpl()      {}

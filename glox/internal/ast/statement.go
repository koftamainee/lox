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

type FunStatement struct {
	Name   token.Token
	Params []token.Token
	Body   []Statement
}

type BlockStatement struct {
	Statements []Statement
}

type IfStatement struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

type WhileStatement struct {
	Condition Expression
	Body      Statement
}

type BreakStatement struct{}

type ReturnStatement struct {
	ReturnKeyword token.Token
	Value         Expression
}

func (e *ExpressionStatement) statementImpl() {}
func (e *PrintStatement) statementImpl()      {}
func (e *VarStatement) statementImpl()        {}
func (e *BlockStatement) statementImpl()      {}
func (e *IfStatement) statementImpl()         {}
func (e *WhileStatement) statementImpl()      {}
func (e *BreakStatement) statementImpl()      {}
func (e *FunStatement) statementImpl()        {}
func (e *ReturnStatement) statementImpl()     {}

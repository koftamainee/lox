package ast

import "github.com/koftamainee/lox/glox/internal/token"

type Expression interface {
	expressionImpl()
}

type BinaryExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

type GroupingExpression struct {
	Expr Expression
}

type LiteralExpression struct {
	Value any
}

type UnaryExpression struct {
	Operator token.Token
	Operand  Expression
}

type ConditionalExpression struct {
	Condition Expression
	Then      Expression
	Else      Expression
}

type VariableExpression struct {
	Name token.Token
}

type AssignmentExpression struct {
	Name  token.Token
	Value Expression
}

type LogicalExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

type CallExpression struct {
	Callee    Expression
	Paren     token.Token
	Arguments []Expression
}

func (e *BinaryExpression) expressionImpl()      {}
func (e *GroupingExpression) expressionImpl()    {}
func (e *LiteralExpression) expressionImpl()     {}
func (e *UnaryExpression) expressionImpl()       {}
func (e *ConditionalExpression) expressionImpl() {}
func (e *VariableExpression) expressionImpl()    {}
func (e *AssignmentExpression) expressionImpl()  {}
func (e *LogicalExpression) expressionImpl()     {}
func (e *CallExpression) expressionImpl()        {}

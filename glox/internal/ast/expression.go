package ast

import (
	"fmt"

	"github.com/koftamainee/lox/glox/internal/token"
)

type Expression interface {
	exprImpl()
	String() string // for not-so-pretty printing
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

func (e *BinaryExpression) exprImpl()      {}
func (e *GroupingExpression) exprImpl()    {}
func (e *LiteralExpression) exprImpl()     {}
func (e *UnaryExpression) exprImpl()       {}
func (e *ConditionalExpression) exprImpl() {}
func (e *VariableExpression) exprImpl()    {}
func (e *AssignmentExpression) exprImpl()  {}
func (e *LogicalExpression) exprImpl()     {}

func (b *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Operator.Lexeme, b.Left.String(), b.Right.String())
}

func (b *LogicalExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Operator.Lexeme, b.Left.String(), b.Right.String())
}

func (g *GroupingExpression) String() string {
	return fmt.Sprintf("(%s)", g.Expr.String())
}

func (l *LiteralExpression) String() string {
	if l.Value == nil {
		return "nil"
	}

	var res string

	switch l.Value.(type) {
	case string:
		res = fmt.Sprintf("\"%v\"", l.Value)
	default:
		res = fmt.Sprintf("%v", l.Value)
	}

	return res
}

func (u *UnaryExpression) String() string {
	return fmt.Sprintf("(%s %s)", u.Operator.Lexeme, u.Operand.String())
}

func (c *ConditionalExpression) String() string {
	return fmt.Sprintf("if %s then %s else %s",
		c.Condition.String(),
		c.Then.String(),
		c.Else.String(),
	)
}

func (v *VariableExpression) String() string {
	return fmt.Sprintf("var %s = %v", v.Name.Lexeme, v.Name.Literal)
}

func (v *AssignmentExpression) String() string {
	return fmt.Sprintf("%s = %v", v.Name.Lexeme, v.Name.Literal)
}

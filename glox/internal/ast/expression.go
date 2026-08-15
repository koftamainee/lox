package ast

import (
	"fmt"

	"github.com/koftamainee/lox/glox/internal/token"
)

type Expression interface {
	exprImpl()
	String() string // for not-so-pretty printing
}

type Binary struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

type Grouping struct {
	Expr Expression
}

type Literal struct {
	Value any
}

type Unary struct {
	Operator token.Token
	Operand  Expression
}

type Conditional struct {
	Condition Expression
	Then      Expression
	Else      Expression
}

func (e *Binary) exprImpl()      {}
func (e *Grouping) exprImpl()    {}
func (e *Literal) exprImpl()     {}
func (e *Unary) exprImpl()       {}
func (e *Conditional) exprImpl() {}

func (b *Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Operator.Lexeme, b.Left.String(), b.Right.String())
}

func (g *Grouping) String() string {
	return fmt.Sprintf("(%s)", g.Expr.String())
}

func (l *Literal) String() string {
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

}

func (u *Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.Operator.Lexeme, u.Operand.String())
}

func (c *Conditional) String() string {
	return fmt.Sprintf("if %s then %s else %s",
		c.Condition.String(),
		c.Then.String(),
		c.Else.String(),
	)
}

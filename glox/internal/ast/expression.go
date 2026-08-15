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

func (b *Binary) exprImpl()   {}
func (b *Grouping) exprImpl() {}
func (b *Literal) exprImpl()  {}
func (b *Unary) exprImpl()    {}

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
	return fmt.Sprintf("%v", l.Value)
}

func (u *Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.Operator.Lexeme, u.Operand.String())
}

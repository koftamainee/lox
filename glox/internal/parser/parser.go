package parser

import (
	"github.com/koftamainee/lox/glox/internal/ast"
	"github.com/koftamainee/lox/glox/internal/token"
)

type Parser struct {
	Tokens []token.Token

	current int
}

func New(tokens []token.Token) Parser {
	return Parser{
		Tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ast.Expression {
	return p.expression()
}

func (p *Parser) match(types ...token.Type) bool {
	for _, token_type := range types {
		if p.check(token_type) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(token_type token.Type) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().TokenType == token_type
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == token.EOF
}

func (p *Parser) peek() *token.Token {
	return &p.Tokens[p.current]
}

func (p *Parser) previous() *token.Token {
	return &p.Tokens[p.current-1]
}

func (p *Parser) expression() ast.Expression {
	return p.equality()
}

func (p *Parser) binaryExprLA(lower func() ast.Expression, operators ...token.Type) ast.Expression {
	expr := lower()

	for p.match(operators...) {
		operator := p.previous()
		right := lower()

		expr = &ast.Binary{
			Left:     expr,
			Operator: *operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) equality() ast.Expression {
	return p.binaryExprLA(p.comparison, token.BangEqual, token.EqualEqual)
}

func (p *Parser) comparison() ast.Expression {
	return p.binaryExprLA(p.term, token.Greater, token.GreaterEqual, token.Less, token.LessEqual)
}

func (p *Parser) term() ast.Expression {
	return p.binaryExprLA(p.factor, token.Plus, token.Minus)
}

func (p *Parser) factor() ast.Expression {
	return p.binaryExprLA(p.unary, token.Star, token.Slash)
}

func (p *Parser) unary() ast.Expression {
	if p.match(token.Bang, token.Equal) {
		operator := p.previous()
		operand := p.unary()
		return &ast.Unary{
			Operator: *operator,
			Operand:  operand,
		}
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expression {
	if p.match(token.True) {
		return &ast.Literal{Value: true}
	}

	if p.match(token.False) {
		return &ast.Literal{Value: false}
	}

	if p.match(token.Nil) {
		return &ast.Literal{Value: nil}
	}

	if p.match(token.Number, token.String) {
		return &ast.Literal{Value: p.previous().Literal}
	}

	if p.match(token.LeftParen) {
		expr := p.expression()
		p.consume(token.RightParen, "Expect ')' after expression.")

		return &ast.Grouping{Expr: expr}
	}

	return nil // Should be unreachable
}

func (p *Parser) consume(token_type token.Type, message string) {
	// TODO
}

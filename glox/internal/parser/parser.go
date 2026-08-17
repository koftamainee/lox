package parser

import (
	"errors"
	"fmt"

	"github.com/koftamainee/lox/glox/internal/ast"
	errrp "github.com/koftamainee/lox/glox/internal/error"
	"github.com/koftamainee/lox/glox/internal/token"
)

var parseErr = errors.New("parse error")

type Parser struct {
	Tokens []token.Token

	errors errrp.ErrorReporter

	current int

	loopDepth int
}

func New(tokens []token.Token, errors errrp.ErrorReporter) Parser {
	return Parser{
		Tokens:    tokens,
		errors:    errors,
		current:   0,
		loopDepth: 0,
	}
}

func (p *Parser) Parse() []ast.Statement {
	statements := make([]ast.Statement, 0)

	for !p.isAtEnd() {
		st := p.declaration()
		if st != nil {
			statements = append(statements, st)
		}
	}

	return statements
}

func (p *Parser) ParseExpression() ast.Expression {
	expr, err := p.expression()
	if err != nil {
		return nil
	}
	if !p.isAtEnd() {
		_ = p.error(p.peek(), "Expect end of input after expression")
		return nil
	}

	return expr
}

func (p *Parser) sync() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == token.Semicolon {
			return
		}

		switch p.peek().TokenType {
		case token.Class:
			return
		case token.Fun:
			return
		case token.Var:
			return
		case token.For:
			return
		case token.If:
			return
		case token.While:
			return
		case token.Print:
			return
		case token.Return:
			return
		case token.Break:
			return
		}

		p.advance()
	}
}

func (p *Parser) binaryOpErrProd(op token.Type) func() (ast.Expression, error) {
	switch op {
	case token.Comma:
		return p.assignment
	case token.Equal:
		return p.ternary
	case token.BangEqual, token.EqualEqual:
		return p.comparison
	case token.Greater, token.GreaterEqual, token.Less, token.LessEqual:
		return p.term
	case token.Plus, token.Minus:
		return p.factor
	case token.Star, token.Slash:
		return p.unary
	}

	return nil
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

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.Tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.Tokens[p.current-1]
}

func (p *Parser) declaration() ast.Statement {
	if p.match(token.Fun) {
		st, err := p.funDeclaration("function")
		if err != nil {
			p.sync()
			return nil
		}
		return st
	}
	if p.match(token.Var) {
		st, err := p.varDeclaration()
		if err != nil {
			p.sync()
			return nil
		}
		return st
	}

	st, err := p.statement()
	if err != nil {
		p.sync()
		return nil
	}
	return st
}

func (p *Parser) funDeclaration(kind string) (ast.Statement, error) {
	name, err := p.consume(token.Identifier, fmt.Sprintf("Expect %s name", kind))
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.LeftParen, fmt.Sprintf("Expect '(' after %s name", kind))
	if err != nil {
		return nil, err
	}

	params := make([]token.Token, 0)
	if !p.check(token.RightParen) {
		for {
			if len(params) >= 255 {
				err := p.error(p.peek(), "Can't have more than 255 parameters")
				if err != nil {
					return nil, err
				}
			}

			param, err := p.consume(token.Identifier, "Expect parameter name")
			if err != nil {
				return nil, err
			}

			params = append(params, param)

			if !p.match(token.Comma) {
				break
			}
		}
	}

	_, err = p.consume(token.RightParen, fmt.Sprintf("Expect ')' after %s parameters", kind))
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.LeftBrace, fmt.Sprintf("Expect '{' before %s body", kind))
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &ast.FunStatement{Name: name, Params: params, Body: body}, nil

}

func (p *Parser) varDeclaration() (ast.Statement, error) {
	name, err := p.consume(token.Identifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expression

	if p.match(token.Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	} else {
		initializer = nil
	}

	_, err = p.consume(token.Semicolon, "Expect semicolon at the end of variable declaration")
	if err != nil {
		return nil, err
	}

	return &ast.VarStatement{
		Name:        name,
		Initializer: initializer,
	}, nil
}

func (p *Parser) statement() (ast.Statement, error) {
	if p.match(token.If) {
		return p.ifStatement()
	}
	if p.match(token.Print) {
		return p.printStatement()
	}
	if p.match(token.While) {
		return p.whileStatement()
	}
	if p.match(token.For) {
		return p.forStatement()
	}
	if p.match(token.Break) {
		return p.breakStatement()
	}
	if p.match(token.Return) {
		return p.returnStatement()
	}

	if p.match(token.LeftBrace) {
		return p.blockStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) returnStatement() (ast.Statement, error) {
	returnKw := p.previous()

	var retval ast.Expression

	if !p.check(token.Semicolon) {
		var err error
		retval, err = p.expression()
		if err != nil {
			return nil, err
		}
	} else {
		retval = nil
	}

	_, err := p.consume(token.Semicolon, "Expect ';' after return statement")
	if err != nil {
		return nil, err
	}

	return &ast.ReturnStatement{Value: retval, ReturnKeyword: returnKw}, nil
}

func (p *Parser) breakStatement() (ast.Statement, error) {
	if p.loopDepth == 0 {
		return nil, p.error(p.previous(), "Expect 'break' statement to be inside the loop")
	}

	_, err := p.consume(token.Semicolon, "Expect ';' after break statement")
	if err != nil {
		return nil, err
	}

	return &ast.BreakStatement{}, nil
}

func (p *Parser) forStatement() (ast.Statement, error) {
	_, err := p.consume(token.LeftParen, "Expect '(' after 'for'")
	if err != nil {
		return nil, err
	}

	var initializer ast.Statement
	if p.match(token.Semicolon) {
		initializer = nil
	} else if p.match(token.Var) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition ast.Expression
	if !p.check(token.Semicolon) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	} else {
		condition = nil
	}
	_, err = p.consume(token.Semicolon, "Expect ';' after condition")
	if err != nil {
		return nil, err
	}

	var increment ast.Expression
	if !p.check(token.LeftParen) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	} else {
		increment = nil
	}

	_, err = p.consume(token.RightParen, "Expect ')' after increment")
	if err != nil {
		return nil, err
	}

	p.loopDepth++
	defer func() { p.loopDepth-- }()

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &ast.BlockStatement{Statements: []ast.Statement{body, &ast.ExpressionStatement{Expr: increment}}}
	}

	if condition == nil {
		condition = &ast.LiteralExpression{Value: true}
	}

	body = &ast.WhileStatement{Condition: condition, Body: body}

	if initializer != nil {
		body = &ast.BlockStatement{Statements: []ast.Statement{initializer, body}}
	}

	return body, nil
}

func (p *Parser) whileStatement() (ast.Statement, error) {
	_, err := p.consume(token.LeftParen, "Expect '(' after 'while'")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RightParen, "Expect ')' after loop condition")
	if err != nil {
		return nil, err
	}

	p.loopDepth++
	defer func() { p.loopDepth-- }()

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStatement{Condition: condition, Body: body}, nil
}

func (p *Parser) ifStatement() (ast.Statement, error) {
	_, err := p.consume(token.LeftParen, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.RightParen, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	var thenBranch ast.Statement
	var elseBranch ast.Statement

	thenBranch, err = p.statement()
	if err != nil {
		return nil, err
	}

	if p.match(token.Else) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	} else {
		elseBranch = nil
	}

	return &ast.IfStatement{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

func (p *Parser) blockStatement() (ast.Statement, error) {
	statements, err := p.block()
	if err != nil {
		return nil, err
	}
	return &ast.BlockStatement{Statements: statements}, nil
}

func (p *Parser) block() ([]ast.Statement, error) {
	statements := make([]ast.Statement, 0)
	for !p.check(token.RightBrace) && !p.isAtEnd() {
		st := p.declaration()
		if st != nil {
			statements = append(statements, st)
		}
	}

	_, err := p.consume(token.RightBrace, "Expect '}' after block")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) printStatement() (ast.Statement, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.Semicolon, "Expect semicolon at the end of print statement.")
	if err != nil {
		return nil, err
	}

	return &ast.PrintStatement{
		Expr: value,
	}, nil
}

func (p *Parser) expressionStatement() (ast.Statement, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.Semicolon, "Expect semicolon at the end of expression statement.")
	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStatement{
		Expr: value,
	}, nil
}

func (p *Parser) expression() (ast.Expression, error) {
	return p.comma()
}

func (p *Parser) comma() (ast.Expression, error) {
	return p.binaryExprLA(p.assignment, token.Comma)
}

func (p *Parser) assignment() (ast.Expression, error) {
	expr, err := p.ternary()
	if err != nil {
		return nil, err
	}

	if p.match(token.Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		e, ok := expr.(*ast.VariableExpression)
		if !ok {
			return nil, p.error(equals, "Variable expected as an assignment target")
		}

		return &ast.AssignmentExpression{Name: e.Name, Value: value}, nil
	}

	return expr, nil
}

func (p *Parser) ternary() (ast.Expression, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(token.Question) {
		thenBranch, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.Colon, "Excpect ':' after then branch of conditional expression.")
		if err != nil {
			return nil, err
		}

		elseBranch, err := p.ternary()
		if err != nil {
			return nil, err
		}

		expr = &ast.ConditionalExpression{
			Condition: expr,
			Then:      thenBranch,
			Else:      elseBranch,
		}
	}

	return expr, nil
}

func (p *Parser) or() (ast.Expression, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(token.Or) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) and() (ast.Expression, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.match(token.And) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) binaryExprLA(
	lower func() (ast.Expression, error),
	operators ...token.Type) (ast.Expression, error) {

	expr, err := lower()
	if err != nil {
		return nil, err
	}

	for p.match(operators...) {
		operator := p.previous()
		right, err := lower()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) equality() (ast.Expression, error) {
	return p.binaryExprLA(p.comparison, token.BangEqual, token.EqualEqual)
}

func (p *Parser) comparison() (ast.Expression, error) {
	return p.binaryExprLA(p.term, token.Greater, token.GreaterEqual, token.Less, token.LessEqual)
}

func (p *Parser) term() (ast.Expression, error) {
	return p.binaryExprLA(p.factor, token.Plus, token.Minus)
}

func (p *Parser) factor() (ast.Expression, error) {
	return p.binaryExprLA(p.unary, token.Star, token.Slash)
}

func (p *Parser) unary() (ast.Expression, error) {
	if p.match(token.Bang, token.Minus) {
		operator := p.previous()
		operand, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpression{
			Operator: operator,
			Operand:  operand,
		}, nil
	}

	return p.call()
}

func (p *Parser) call() (ast.Expression, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(token.LeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee ast.Expression) (ast.Expression, error) {
	args := make([]ast.Expression, 0)

	if !p.check(token.RightParen) {
		for {
			if len(args) >= 255 {
				// NOTE(koftamainee): do not return error, parser still in valid state
				_ = p.error(p.peek(), "Can't have more than 255 arguments")
			}
			// NOTE(koftamainee): skipping comma operator here
			expr, err := p.assignment()
			if err != nil {
				return nil, err
			}
			args = append(args, expr)

			if !p.match(token.Comma) {
				break
			}
		}
	}
	paren, err := p.consume(token.RightParen, "Expect ')' after arguments")
	if err != nil {
		return nil, err
	}

	return &ast.CallExpression{Callee: callee, Paren: paren, Arguments: args}, nil
}

func (p *Parser) primary() (ast.Expression, error) {
	if p.match(token.True) {
		return &ast.LiteralExpression{Value: true}, nil
	}

	if p.match(token.False) {
		return &ast.LiteralExpression{Value: false}, nil
	}

	if p.match(token.Nil) {
		return &ast.LiteralExpression{Value: nil}, nil
	}

	if p.match(token.Number, token.String) {
		return &ast.LiteralExpression{Value: p.previous().Literal}, nil
	}

	if p.match(token.Identifier) {
		return &ast.VariableExpression{Name: p.previous()}, nil
	}

	if p.match(token.LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return &ast.GroupingExpression{Expr: expr}, nil
	}

	operand := p.binaryOpErrProd(p.peek().TokenType)
	if operand != nil {
		p.advance()
		p.error(p.previous(), "Expect expression.")
		return operand()

	}

	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) consume(token_type token.Type, message string) (token.Token, error) {
	if p.check(token_type) {
		return p.advance(), nil
	}

	return token.Token{}, p.error(p.peek(), message)
}

func (p *Parser) error(token token.Token, message string) error {
	p.errors.ErrorAt(token, message)
	return fmt.Errorf("%w: %s", parseErr, message)
}

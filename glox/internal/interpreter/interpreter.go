package interpreter

import (
	"errors"
	"fmt"

	"github.com/koftamainee/lox/glox/internal/ast"
	errrp "github.com/koftamainee/lox/glox/internal/error"
	"github.com/koftamainee/lox/glox/internal/token"
)

type RuntimeError struct {
	Token token.Token
	Msg   string
}

func (e RuntimeError) Error() string {
	return e.Msg
}

type Interpreter struct {
	errors errrp.ErrorReporter
}

func New(errors errrp.ErrorReporter) Interpreter {
	return Interpreter{
		errors: errors,
	}
}

func (i *Interpreter) Interpret(expr ast.Expression) any {
	value, err := i.evaluateExpression(expr)
	if err != nil {
		runtimeError, ok := errors.AsType[RuntimeError](err)
		if ok {
			i.errors.RuntimeError(runtimeError.Token, runtimeError.Msg)
		} else {
			i.errors.InternalError(err.Error())
		}
		return nil
	}

	return value
}

func (i *Interpreter) evaluateExpression(expr ast.Expression) (any, error) {
	switch e := expr.(type) {
	case *ast.Literal:
		return i.evalLiteralExpr(e)
	case *ast.Grouping:
		return i.evalGroupingExpr(e)
	case *ast.Binary:
		return i.evalBinaryExpr(e)
	case *ast.Unary:
		return i.evalUnaryExpr(e)
	case *ast.Conditional:
		return i.evalConditionalExpr(e)
	}

	return nil, errors.New("invalid expression type")
}

func (i *Interpreter) evalLiteralExpr(expr *ast.Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) evalGroupingExpr(expr *ast.Grouping) (any, error) {
	return i.evaluateExpression(expr.Expr)
}

func (i *Interpreter) evalBinaryExpr(expr *ast.Binary) (any, error) {
	left, err := i.evaluateExpression(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluateExpression(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case token.Minus:
		leftv, err := checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv - rightv, nil

	case token.Star:
		leftv, err := checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv * rightv, nil

	case token.Slash:
		leftv, err := checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv / rightv, nil

	case token.Plus:
		leftv, ok := left.(float64)
		if ok {
			rightv, ok := right.(float64)
			if !ok {
				return nil, rtError(expr.Operator, "Operands must be two numbers or two strings")
			}
			return leftv + rightv, nil
		}
		lefts, ok := left.(string)
		if ok {
			rights, ok := right.(string)
			if !ok {
				return nil, rtError(expr.Operator, "Operands must be two numbers or two strings")
			}

			return fmt.Sprintf("%s%s", lefts, rights), nil
		} else {
			return nil, rtError(expr.Operator, "Operands must be two numbers or two strings")
		}

	case token.Greater:
		leftv, err := checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv > rightv, nil

	case token.GreaterEqual:
		leftv, err := checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv >= rightv, nil

	case token.Less:
		leftv, err := checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv < rightv, nil

	case token.LessEqual:
		leftv, err := checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv <= rightv, nil

	case token.EqualEqual:
		return isEqual(left, right), nil

	case token.BangEqual:
		return !isEqual(left, right), nil

	case token.Comma:
		return right, nil
	default:
		return nil, rtError(expr.Operator, fmt.Sprintf("Invalid operator: '%s'", expr.Operator.Lexeme))
	}

}

func (i *Interpreter) evalUnaryExpr(expr *ast.Unary) (any, error) {
	value, err := i.evaluateExpression(expr.Operand)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case token.Minus:
		switch v := value.(type) {
		case float64:
			return -v, nil
		default:
			return nil, rtError(expr.Operator, fmt.Sprintf("Invalid operator '%s' on value '%v'", expr.Operator.Lexeme, v))
		}
	case token.Bang:
		// NOTE(koftamainee): every value other then nil and false are true. this is questionable desision
		return !isTruthy(value), nil
	default:
		return nil, rtError(expr.Operator, fmt.Sprintf("Invalid operator: '%s'", expr.Operator.Lexeme))
	}
}

func (i *Interpreter) evalConditionalExpr(expr *ast.Conditional) (any, error) {
	condition, err := i.evaluateExpression(expr.Condition)
	if err != nil {
		return nil, err
	}
	if isTruthy(condition) {
		return i.evaluateExpression(expr.Then)
	} else {
		return i.evaluateExpression(expr.Else)
	}
}

func rtError(token token.Token, msg string) RuntimeError {
	return RuntimeError{
		Token: token,
		Msg:   msg,
	}
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case bool:
		return v
	default:
		return true
	}
}

func isEqual(left any, right any) bool {
	if left == nil {
		return right == nil
	}

	return left == right
}

func checkNumber(operator token.Token, operand any) (float64, error) {
	switch v := operand.(type) {
	case float64:
		return v, nil
	default:
		return 0, rtError(operator, "Operand must be a number")
	}
}

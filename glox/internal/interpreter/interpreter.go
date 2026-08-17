package interpreter

import (
	"errors"
	"fmt"

	"github.com/koftamainee/lox/glox/internal/ast"
	"github.com/koftamainee/lox/glox/internal/environment"
	errrp "github.com/koftamainee/lox/glox/internal/error"
	"github.com/koftamainee/lox/glox/internal/token"
)

var breakSt = errors.New("break")

type RuntimeError struct {
	Token token.Token
	Msg   string
}

func (e RuntimeError) Error() string {
	return e.Msg
}

type Interpreter struct {
	errors errrp.ErrorReporter

	globals    *environment.Environment
	currentEnv *environment.Environment
}

func New(errors errrp.ErrorReporter) Interpreter {
	globals := environment.New(nil)
	defineGlobals(globals)
	return Interpreter{
		errors:     errors,
		globals:    globals,
		currentEnv: globals,
	}
}

func (i *Interpreter) Interpret(statements []ast.Statement) {

	for _, st := range statements {
		err := i.executeStatement(st)
		if err != nil {
			runtimeError, ok := errors.AsType[RuntimeError](err)
			if ok {
				i.errors.RuntimeError(runtimeError.Token, runtimeError.Msg)
			} else {
				i.errors.InternalError(err.Error())
			}
		}
	}
}

func (i *Interpreter) Evaluate(expr ast.Expression) any {
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

func (i *Interpreter) executeStatement(st ast.Statement) error {
	switch s := st.(type) {
	case *ast.ExpressionStatement:
		return i.execExpressionStmt(s)
	case *ast.PrintStatement:
		return i.execPrintStmt(s)
	case *ast.VarStatement:
		return i.execVarStmt(s)
	case *ast.BlockStatement:
		return i.execBlockStmt(s, environment.New(i.currentEnv))
	case *ast.IfStatement:
		return i.execIfStmt(s)
	case *ast.WhileStatement:
		return i.execWhileStmt(s)
	case *ast.BreakStatement:
		return i.execBreakStmt(s)
	case *ast.FunStatement:
		return i.execFunStmt(s)
	case *ast.ReturnStatement:
		return i.execReturnStmt(s)

	default:
		return errors.New("invalid statement type")
	}
}

func (i *Interpreter) execReturnStmt(st *ast.ReturnStatement) error {
	var retval any

	if st.Value != nil {
		var err error
		retval, err = i.evaluateExpression(st.Value)
		if err != nil {
			return err
		}
	} else {
		retval = nil
	}

	return returnErr{Value: retval}
}

func (i *Interpreter) execFunStmt(st *ast.FunStatement) error {
	function := loxFunction{Declaration: st}
	i.currentEnv.Define(st.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) execBreakStmt(st *ast.BreakStatement) error {
	_ = st
	return breakSt
}

func (i *Interpreter) execWhileStmt(st *ast.WhileStatement) error {

	for {
		condition, err := i.evaluateExpression(st.Condition)
		if err != nil {
			return err
		}
		if !isTruthy(condition) {
			break
		}

		err = i.executeStatement(st.Body)
		if err != nil {
			if errors.Is(err, breakSt) {
				break
			}
			return err
		}
	}

	return nil
}

func (i *Interpreter) execIfStmt(st *ast.IfStatement) error {
	cond, err := i.evaluateExpression(st.Condition)
	if err != nil {
		return err
	}

	if isTruthy(cond) {
		err := i.executeStatement(st.ThenBranch)
		if err != nil {
			return err
		}
	} else if st.ElseBranch != nil {
		err := i.executeStatement(st.ElseBranch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) execBlockStmt(st *ast.BlockStatement, env *environment.Environment) error {
	previousEnv := i.currentEnv

	i.currentEnv = env
	defer func() { i.currentEnv = previousEnv }()

	for _, st := range st.Statements {
		err := i.executeStatement(st)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) execExpressionStmt(st *ast.ExpressionStatement) error {
	_, err := i.evaluateExpression(st.Expr)
	return err
}

func (i *Interpreter) execPrintStmt(st *ast.PrintStatement) error {
	value, err := i.evaluateExpression(st.Expr)
	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

func (i *Interpreter) execVarStmt(st *ast.VarStatement) error {
	if st.Initializer != nil {
		value, err := i.evaluateExpression(st.Initializer)
		if err != nil {
			return err
		}
		i.currentEnv.Define(st.Name.Lexeme, value)
	} else {
		i.currentEnv.Declare(st.Name.Lexeme)
	}

	return nil
}

func (i *Interpreter) evaluateExpression(expr ast.Expression) (any, error) {
	switch e := expr.(type) {
	case *ast.LiteralExpression:
		return i.evalLiteralExpr(e)
	case *ast.GroupingExpression:
		return i.evalGroupingExpr(e)
	case *ast.BinaryExpression:
		return i.evalBinaryExpr(e)
	case *ast.UnaryExpression:
		return i.evalUnaryExpr(e)
	case *ast.ConditionalExpression:
		return i.evalConditionalExpr(e)
	case *ast.VariableExpression:
		return i.evalVariableExpr(e)
	case *ast.AssignmentExpression:
		return i.evalAssignmentExpr(e)
	case *ast.LogicalExpression:
		return i.evalLogicalExpr(e)
	case *ast.CallExpression:
		return i.evalCallExpr(e)

	default:
		return nil, errors.New("invalid expression type")
	}
}

func (i *Interpreter) evalCallExpr(expr *ast.CallExpression) (any, error) {
	callee, err := i.evaluateExpression(expr.Callee)
	if err != nil {
		return nil, err
	}

	args := make([]any, len(expr.Arguments))
	for index, arg := range expr.Arguments {
		argv, err := i.evaluateExpression(arg)
		if err != nil {
			return nil, err
		}
		args[index] = argv
	}

	callable, ok := callee.(loxCallable)
	if !ok {
		return nil, i.error(expr.Paren, "Expected Object to be LoxCallable")
	}

	expectedArgs := callable.Arity()
	gotArgs := len(args)

	if expectedArgs != gotArgs {
		return nil, i.error(expr.Paren, fmt.Sprintf("Expected %d arguments but got %d", expectedArgs, gotArgs))
	}

	return callable.Call(i, args)
}

func (i *Interpreter) evalLogicalExpr(expr *ast.LogicalExpression) (any, error) {
	left, err := i.evaluateExpression(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.TokenType == token.Or {
		if isTruthy(left) {
			return left, nil
		}
	} else { // token.And
		if !isTruthy(left) {
			return left, nil
		}
	}
	return i.evaluateExpression(expr.Right)
}

func (i *Interpreter) evalAssignmentExpr(expr *ast.AssignmentExpression) (any, error) {
	value, err := i.evaluateExpression(expr.Value)
	if err != nil {
		return nil, err
	}
	err = i.currentEnv.Assign(expr.Name, value)
	if err != nil {
		if errors.Is(err, environment.ErrUndeclared) {
			return nil, i.error(expr.Name, "Variable is undeclared")
		}
		return nil, err
	}

	return value, nil
}

func (i *Interpreter) evalVariableExpr(expr *ast.VariableExpression) (any, error) {
	value, err := i.currentEnv.Get(expr.Name)
	if err != nil {
		if errors.Is(err, environment.ErrUndefined) {
			return nil, i.error(expr.Name, "Accessing undefined variable. Please assign value to it")
		} else if errors.Is(err, environment.ErrUndeclared) {
			return nil, i.error(expr.Name, "Accessing undeclared variable.")
		}
		// NOTE(koftamainee): propagating all possible internal errors
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) evalLiteralExpr(expr *ast.LiteralExpression) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) evalGroupingExpr(expr *ast.GroupingExpression) (any, error) {
	return i.evaluateExpression(expr.Expr)
}

func (i *Interpreter) evalBinaryExpr(expr *ast.BinaryExpression) (any, error) {
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
		leftv, err := i.checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := i.checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv - rightv, nil

	case token.Star:
		leftv, err := i.checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := i.checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv * rightv, nil

	case token.Slash:
		leftv, err := i.checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := i.checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv / rightv, nil

	case token.Plus:
		leftv, ok := left.(float64)
		if ok {
			rightv, ok := right.(float64)
			if !ok {
				return nil, i.error(expr.Operator, "Operands must be two numbers or two strings")
			}
			return leftv + rightv, nil
		}
		lefts, ok := left.(string)
		if ok {
			rights, ok := right.(string)
			if !ok {
				return nil, i.error(expr.Operator, "Operands must be two numbers or two strings")
			}

			return fmt.Sprintf("%s%s", lefts, rights), nil
		} else {
			return nil, i.error(expr.Operator, "Operands must be two numbers or two strings")
		}

	case token.Greater:
		leftv, err := i.checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := i.checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv > rightv, nil

	case token.GreaterEqual:
		leftv, err := i.checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := i.checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv >= rightv, nil

	case token.Less:
		leftv, err := i.checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := i.checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return leftv < rightv, nil

	case token.LessEqual:
		leftv, err := i.checkNumber(expr.Operator, left)
		if err != nil {
			return nil, err
		}
		rightv, err := i.checkNumber(expr.Operator, right)
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
		return nil, i.error(expr.Operator, fmt.Sprintf("Invalid operator: '%s'", expr.Operator.Lexeme))
	}

}

func (i *Interpreter) evalUnaryExpr(expr *ast.UnaryExpression) (any, error) {
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
			return nil, i.error(expr.Operator, fmt.Sprintf("Invalid operator '%s' on value '%v'", expr.Operator.Lexeme, v))
		}
	case token.Bang:
		// NOTE(koftamainee): every value other then nil and false are true. this is questionable desision
		return !isTruthy(value), nil
	default:
		return nil, i.error(expr.Operator, fmt.Sprintf("Invalid operator: '%s'", expr.Operator.Lexeme))
	}
}

func (i *Interpreter) evalConditionalExpr(expr *ast.ConditionalExpression) (any, error) {
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

func (i *Interpreter) error(token token.Token, msg string) RuntimeError {
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

func (i *Interpreter) checkNumber(operator token.Token, operand any) (float64, error) {
	switch v := operand.(type) {
	case float64:
		return v, nil
	default:
		return 0, i.error(operator, "Operand must be a number")
	}
}

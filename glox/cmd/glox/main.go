package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/koftamainee/lox/glox/internal/ast"
	"github.com/koftamainee/lox/glox/internal/lox"
	"github.com/koftamainee/lox/glox/internal/token"
)

const (
	Success           = 0
	WrongArgsCountErr = 64
	DataFormatErr     = 65
	GLoxInternalErr   = 70
)

func main() {

	minus_123 := ast.Unary{
		Operator: token.New(token.Minus, "-", nil, 1),
		Operand:  &ast.Literal{Value: 123},
	}

	group45_67 := ast.Grouping{
		Expr: &ast.Literal{Value: 45.67},
	}

	expr := ast.Binary{
		Left:     &minus_123,
		Operator: token.New(token.Star, "*", nil, 1),
		Right:    &group45_67,
	}

	fmt.Printf("%s\n", expr.String())

	if len(os.Args) > 2 {
		fmt.Println("Usage glox [script]")
		os.Exit(WrongArgsCountErr)
	}

	lox, err := lox.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "glox encountered fatal error during initialization: %e\n", err)
		os.Exit(GLoxInternalErr)
	}

	if len(os.Args) == 2 {
		runFile(os.Args[1], &lox)
	} else {
		runPrompt(&lox)
	}

	if lox.ErrorReporter.HadError {
		os.Exit(DataFormatErr)
	}

	os.Exit(Success)
}

func runFile(path string, lox *lox.Lox) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		lox.ErrorReporter.InternalError(fmt.Errorf("failed to open file %s: %v", path, err).Error())
	}
	lox.Run(string(bytes))
}

func runPrompt(lox *lox.Lox) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("glox v0.1.0")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			fmt.Println("\nsee you soon~")
			break
		}
		line := scanner.Text()
		lox.Run(line)
		lox.ErrorReporter.HadError = false
	}
}

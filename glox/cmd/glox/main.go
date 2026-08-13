package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/koftamainee/lox/glox/internal/lox"
)

const (
	WrongArgsCountErr = 64
	GLoxInternalErr   = 70
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage glox [script]")
		os.Exit(WrongArgsCountErr)
	}

	lox, err := lox.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "glox encountered fatal error during initialization: %e\n", err)
		os.Exit(GLoxInternalErr)
	}

	if len(os.Args) == 2 {
		runFile(os.Args[1], &lox)
	} else {
		runPrompt(&lox)
	}
}

func runFile(path string, lox *lox.Lox) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		lox.InternalError(fmt.Errorf("failed to open file %s: %v", path, err).Error())
	}
	lox.Run(bytes)
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
		line := scanner.Bytes()
		lox.Run(line)
	}
}

func run(bytes []byte) error {
	return nil
}

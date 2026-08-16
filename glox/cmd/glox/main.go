package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/koftamainee/lox/glox/internal/lox"
)

const (
	Success           = 0
	WrongArgsCountErr = 64
	DataFormatErr     = 65
	RuntimeErr        = 70
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage glox [script]")
		os.Exit(WrongArgsCountErr)
	}

	lox := lox.New()
	if len(os.Args) == 2 {
		runFile(os.Args[1], &lox)
	} else {
		runPrompt(&lox)
	}

	if lox.ErrorReporter.HadError {
		os.Exit(DataFormatErr)
	}

	if lox.ErrorReporter.HadRuntimeError {
		os.Exit(RuntimeErr)

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
	fmt.Println("glox v0.1.0")

	rl, err := readline.New("> ")
	if err != nil {
		lox.ErrorReporter.InternalError("failed to create readline")
		return
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			fmt.Println("see you soon~")
			break
		}
		lox.Run(line)
		lox.ErrorReporter.HadError = false
		lox.ErrorReporter.HadRuntimeError = false
	}
}

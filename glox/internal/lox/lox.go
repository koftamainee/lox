package lox

import (
	"fmt"
	"os"
)

type Lox struct {
}

func Init() (Lox, error) {
	return Lox{}, nil
}

func (l *Lox) Run(bytes []byte) error {
	// ```java
	// Scanner scanner = new Scanner(source);
	// List<Token> tokens = scanner.scanTokens();

	// // For now, just print the tokens.
	// for (Token token : tokens) {
	//   System.out.println(token);
	// }
	//```
	return nil
}

func (l *Lox) Error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) InternalError(message string) {
	fmt.Fprintf(os.Stderr, "Glox internal error: %s", message)
}

func (l *Lox) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
}

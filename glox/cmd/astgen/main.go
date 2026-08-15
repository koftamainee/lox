package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <output_directory>", os.Args[0])
		os.Exit(64)
	}

	outputDir := os.Args[1]
	defineAst()

}

func defineAst(outputDir string, baseName string, types []string) {

}

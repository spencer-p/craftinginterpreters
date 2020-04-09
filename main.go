package main

import (
	"flag"

	"github.com/spencer-p/craftinginterpreters/pkg/lox"
)

func main() {
	// TODO Flags?
	flag.Parse()

	inputFile := flag.Arg(0)
	if inputFile == "" {
		lox.RunPrompt()
	} else {
		lox.RunFile(inputFile)
	}
}

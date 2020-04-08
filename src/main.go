package main

import (
	"flag"
)

func main() {
	// TODO Flags?
	flag.Parse()

	inputFile := flag.Arg(0)
	if inputFile == "" {
		runPrompt()
	} else {
		runFile(inputFile)
	}
}

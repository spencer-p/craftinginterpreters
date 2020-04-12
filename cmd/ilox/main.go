package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spencer-p/craftinginterpreters/pkg/lox"
)

func main() {
	// TODO Flags?
	flag.Parse()
	inputFile := flag.Arg(0)

	var err error
	if inputFile == "" {
		err = lox.RunPrompt()
	} else {
		err = lox.RunFile(inputFile)
	}

	if err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", err)
	}
}

package lox

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/chzyer/readline"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/parse"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/prettyprint"
	"github.com/spencer-p/craftinginterpreters/pkg/lox/scan"
)

// RunFile interprets the code in the given file.
func RunFile(path string) error {
	bytes, err := fetchFile(path)
	if err != nil {
		return err
	}

	// free utf-8 support! thanks, go
	run(string(bytes))
	return nil
}

// RunPrompt interprets code interactively.
func RunPrompt() error {
	rl, err := readline.New("> ")
	if err != nil {
		return fmt.Errorf("could not run interactive: %v", err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return fmt.Errorf("failed to read user input: %v", err)
			}
		}
		run(line)
	}
	return nil
}

// fetchFile turns a path into the bytes of the corresponding file.
// The file is loaded into memory and the file resource is cleaned up.
func fetchFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func run(in string) {
	toks := scan.New(in).Tokens()
	ast, err := parse.New(toks).AST()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Println(ast.Accept(&prettyprint.Lisp{}))
}

package lox

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spencer-p/craftinginterpreters/pkg/lox/scan"
)

// runFile interprets the code in the given file.
func RunFile(path string) error {
	bytes, err := fetchFile(path)
	if err != nil {
		return err
	}

	// free utf-8 support! thanks, go
	run(string(bytes))
	return nil
}

// runPrompt interprets code interactively.
func RunPrompt() {
	in := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")
	for in.Scan() {
		run(in.Text())
		fmt.Print("> ")
	}
	fmt.Println()
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
	scanner := scan.New(in)
	toks := scanner.Tokens()
	fmt.Printf("%v\n", toks)
}

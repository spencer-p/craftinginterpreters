/*
This tool generates expressions.

Usage (indented to prevent it running for this example):
	//go:generate go run github.com/spencer-p/craftinginterpreters/cmd/genexpr
	/// Expr1: x int, y string
	/// Expr2: a bool, b []byte

It is also a good idea to format afterwards.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spencer-p/craftinginterpreters/pkg/tool/meta"
)

const (
	PREFIX      = "/// "
	FNAME_KEY   = "GOFILE"
	PACKAGE_KEY = "GOPACKAGE"
)

func die(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	filename := os.Getenv(FNAME_KEY)
	packagename := os.Getenv(PACKAGE_KEY)
	if filename == "" || packagename == "" {
		fmt.Fprintf(os.Stderr, "missing environment vars\n")
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open file\n")
		return
	}
	defer f.Close()

	types := make([]meta.Typ, 0)

	lines := bufio.NewScanner(f)
	for lines.Scan() {
		text := lines.Text()
		if !strings.HasPrefix(text, PREFIX) {
			continue
		}

		text = text[len(PREFIX):]
		parsed, err := meta.ParseTypes(text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse\n")
			return
		}
		types = append(types, parsed...)
	}

	out, err := os.Create(filename[:len(filename)-len(".go")] + "_genexpr.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot create generated file")
		return
	}
	defer f.Close()

	meta.GenExpr(out, types, packagename)
}

package main

import (
	"fmt"
	"os"
)

func complain(line int, msg string) {
	report(line, "", msg)
}

func report(line int, where string, msg string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, msg)
}

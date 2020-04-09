package main

import (
	"fmt"
	"os"
)

func complain(line int, msg string, extra ...interface{}) {
	report(line, "", fmt.Sprintf(msg, extra...))
}

func report(line int, where string, msg string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, msg)
}

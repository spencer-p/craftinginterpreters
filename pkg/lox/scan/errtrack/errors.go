package errtrack

import (
	"fmt"
	"io"
	"os"
)

var (
	hadError           = false
	output   io.Writer = os.Stdout
)

func SetOutput(newOutput io.Writer) {
	output = newOutput
}

func Complain(line, char int, msg string, extra ...interface{}) {
	Report(line, char, "", fmt.Sprintf(msg, extra...))
}

func Report(line, char int, where string, msg string) {
	hadError = true
	fmt.Fprintf(output, "[line %d:%d] Error%s: %s\n", line, char, where, msg)
}

func Err() bool {
	return hadError
}

func Reset() {
	hadError = false
}

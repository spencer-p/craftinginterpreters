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

func Complain(line int, msg string, extra ...interface{}) {
	Report(line, "", fmt.Sprintf(msg, extra...))
}

func Report(line int, where string, msg string) {
	hadError = true
	fmt.Fprintf(output, "[line %d] Error%s: %s\n", line, where, msg)
}

func Err() bool {
	return hadError
}

func Reset() {
	hadError = false
}

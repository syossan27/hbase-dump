package common

import (
	"fmt"
	"os"
	"github.com/mattn/go-colorable"
)

const (
	// ExitCodeOK is exit code for OK
	ExitCodeOK = iota
	// ExitCodeError is exit code for Error
	ExitCodeError
)

var (
	stdout = colorable.NewColorableStdout()
	stderr = colorable.NewColorableStderr()
)

func Success(msg string) {
	fmt.Fprintf(stdout, "\x1b[32m%s\x1b[0m\n", msg)
}

func Fatal(msg string) {
	fmt.Fprintf(stderr, "\x1b[31m%s\x1b[0m\n", msg)
	os.Exit(ExitCodeError)
}

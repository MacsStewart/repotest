package builtins

import (
	"fmt"
	"io"
	"strings"
)

func echoCommand(w io.Writer, args ...string) error {
	// Join the arguments with a space and print them
	echoString := strings.Join(args, " ")
	fmt.Fprintln(w, echoString)
	return nil
}

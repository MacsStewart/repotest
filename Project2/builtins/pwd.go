package builtins

import (
	"fmt"
	"io"
	"os"
)

func printWorkingDirectory(w io.Writer, args ...string) error {
	// Check for valid options
	for _, arg := range args {
		if arg != "-P" && arg != "-L" {
			return fmt.Errorf("pwd: invalid option: %s", arg)
		}
	}

	// Get current working directory
	var (
		wd  string
		err error
	)

	// Check for options
	if len(args) > 0 && args[0] == "-P" {
		wd, err = os.Getwd()
	} else {
		wd, err = os.Getwd()
	}

	if err != nil {
		return fmt.Errorf("pwd: %v", err)
	}

	// Print the current working directory
	fmt.Fprintln(w, wd)
	return nil
}

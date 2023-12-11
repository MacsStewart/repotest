package builtins

import (
	"fmt"
	"io"
	"io/ioutil"
)

func listFiles(w io.Writer, args ...string) error {
	// Use the current directory if no path is provided
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	// Print the list of files
	for _, info := range fileInfo {
		fmt.Fprintf(w, "%s\t", info.Name())
	}
	fmt.Fprintln(w) // Add a newline after the file list

	return nil
}

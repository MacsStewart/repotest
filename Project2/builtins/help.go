package builtins

import (
	"fmt"
	"io"
)

func showHelp(w io.Writer) error {
	helpMessage := `
 Available commands:
 - cd [directory]: Change the current working directory
 - env: Display environment variables
 - caller: Display caller information
 - ls [directory]: List files in the specified directory (default is current directory)
 - pwd [-P | -L]: Print current working directory (-P: physical, -L: logical, default is logical)
 - echo [text]: Display the provided text
 - exit: Exit the shell
 
 Type 'help [command]' for more information on a specific command.
 `

	_, err := fmt.Fprint(w, helpMessage)
	return err
}

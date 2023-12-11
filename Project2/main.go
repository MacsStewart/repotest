package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/MacsStewart/repotest/Project2/builtins"
)

func main() {
	exit := make(chan struct{}, 2) // buffer this so there's no deadlock.
	runLoop(os.Stdin, os.Stdout, os.Stderr, exit)
}

func runLoop(r io.Reader, w, errW io.Writer, exit chan struct{}) {
	var (
		input    string
		err      error
		readLoop = bufio.NewReader(r)
	)
	for {
		select {
		case <-exit:
			_, _ = fmt.Fprintln(w, "exiting gracefully...")
			return
		default:
			if err := printPrompt(w); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if input, err = readLoop.ReadString('\n'); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if err = handleInput(w, input, exit); err != nil {
				_, _ = fmt.Fprintln(errW, err)
			}
		}
	}
}

func printPrompt(w io.Writer) error {
	// Get current user.
	// Don't prematurely memoize this because it might change due to `su`?
	u, err := user.Current()
	if err != nil {
		return err
	}
	// Get current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// /home/User [Username] $
	_, err = fmt.Fprintf(w, "%v [%v] $ ", wd, u.Username)

	return err
}

func handleInput(w io.Writer, input string, exit chan<- struct{}) error {
	// Remove trailing spaces.
	input = strings.TrimSpace(input)

	// Split the input separate the command name and the command arguments.
	args := strings.Split(input, " ")
	name, args := args[0], args[1:]

	// Check for built-in commands.
	// New builtin commands should be added here. Eventually this should be refactored to its own func.
	switch name {
	case "cd":
		return builtins.ChangeDirectory(args...)
	case "env":
		return builtins.EnvironmentVariables(w, args...)
	case "caller":
		return printCallerInfo(w)
	case "ls":
		return listFiles(w, args...)
	case "pwd":
		return printWorkingDirectory(w, args...)
	case "echo":
		return echoCommand(w, args...)
	case "help":
		return showHelp(w)
	case "exit":
		exit <- struct{}{}
		return nil
	}

	return executeCommand(name, args...)
}

func executeCommand(name string, arg ...string) error {
	// Otherwise prep the command
	cmd := exec.Command(name, arg...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}

// Repo was not working This is what will have to work for now

// Caller Function start:
func printCallerInfo(w io.Writer) error {
	callStackInfo := getCallStackInfo(1)
	_, err := fmt.Fprintln(w, callStackInfo)
	return err
}

func getCallStackInfo(levels int) string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	for i := 0; i < levels; i++ {
		frame, more := frames.Next()
		if !more {
			break
		}

		if i == levels-1 {
			return fmt.Sprintf("File: %s, Line: %d, Function: %s", frame.File, frame.Line, frame.Function)
		}
	}

	return "Unable to retrieve call stack information."
}

// caller function end

// ls Function start
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

//ls end

// pwd start
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

// end of pwd
// Start of echo command
func echoCommand(w io.Writer, args ...string) error {
	// Join the arguments with a space and print them
	echoString := strings.Join(args, " ")
	fmt.Fprintln(w, echoString)
	return nil
}

// end of echo command
// Start of Help command
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

//End of help command

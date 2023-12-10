package builtins

import (
	"fmt"
	"runtime"
)

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

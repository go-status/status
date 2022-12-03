package stacktrace

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	stpb "github.com/go-status/status/proto"
)

// ToString returns a human-readable string representing the stack trace.
// The format may change, so callers should not rely on it. If the `verbose`
// flag is enabled, a stack frame will include a fully qualified function
// name and a full path (e.g., "path/to/pkg.Func@/path/to/format.go:123").
// Otherwise, a stack frame will include short names (e.g.,
// "Func@format.go:123").
func ToString(s *stpb.StackTrace, verbose bool) string {
	// If there are no frames, return "No stack trace" for better readability.
	if len(s.GetFrames()) == 0 {
		return "No stack trace"
	}

	// Use a string builder to concatenate strings efficiently.
	var w strings.Builder

	// Append a stack trace heading.
	w.WriteString("Stack trace:")

	// Append all frames in a human-readable format.
	for _, frame := range s.GetFrames() {
		// Indent each stack frame by two spaces.
		w.WriteString("\n  ")

		// Append the function name.
		if verbose {
			w.WriteString(frame.GetFunction())
		} else {
			p := strings.LastIndex(frame.GetFunction(), "/")
			w.WriteString(frame.GetFunction()[p+1:])
		}

		// Append a separator between the function name and the file name.
		w.WriteString("@")

		// Append the file name (verbose=false) or the full path
		// (verbose=true).
		if verbose {
			w.WriteString(frame.GetFile())
		} else {
			w.WriteString(path.Base(frame.GetFile()))
		}

		// Append the line number.
		w.WriteString(":")
		w.WriteString(strconv.Itoa(int(frame.GetLine())))

		// Append the program counter if verbose=true and program counter is
		// not 0.
		if verbose && frame.GetProgramCounter() != 0 {
			w.WriteString("(")

			if frame.GetEntry() != 0 {
				w.WriteString(fmt.Sprintf("0x%x+", frame.GetEntry()))
			}

			w.WriteString(fmt.Sprintf(
				"0x%x)", frame.GetProgramCounter()-frame.GetEntry()))
		}
	}

	// Build the string and return it.
	return w.String()
}

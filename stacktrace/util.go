package stacktrace

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	stpb "github.com/go-status/status/proto"
)

// ToString returns a string describing a stack trace in a human-readable
// format.  The format may change, so callers must not depend on it.  If
// verbose is enabled, a stack frame contains a fully qualified function name
// and a full path (e.g., "path/to/pkg.Func@/path/to/format.go:123").
// Otherwise, a stack frame contains short names (e.g., "Func@format.go:123").
func ToString(s *stpb.StackTrace, verbose bool) string {
	// If no frames exist, returns "No stack trace" for better readability.
	if len(s.GetFrames()) == 0 {
		return "No stack trace"
	}

	// Prepare a string builder because concatenating strings is inefficient.
	var w strings.Builder

	// Append a stack trace heading.
	w.WriteString("Stack trace:")

	// Append all frames in a human-readable format.
	for _, frame := range s.GetFrames() {
		// Make every stack frame have an indent of two spaces.
		w.WriteString("\n  ")

		// Append a function name.
		if verbose {
			w.WriteString(frame.GetFunction())
		} else {
			p := strings.LastIndex(frame.GetFunction(), "/")
			w.WriteString(frame.GetFunction()[p+1:])
		}

		// Append a separator between a function name and a file name.
		w.WriteString("@")

		// Append a file name (verbose=false) or a full path (verbose=false).
		if verbose {
			w.WriteString(frame.GetFile())
		} else {
			w.WriteString(path.Base(frame.GetFile()))
		}

		// Append a line number.
		w.WriteString(":")
		w.WriteString(strconv.Itoa(int(frame.GetLine())))

		// Append a program counter.
		if verbose && frame.GetProgramCounter() != 0 {
			w.WriteString("(")

			if frame.GetEntry() != 0 {
				w.WriteString(fmt.Sprintf("0x%x+", frame.GetEntry()))
			}

			w.WriteString(fmt.Sprintf(
				"0x%x)", frame.GetProgramCounter()-frame.GetEntry()))
		}
	}

	// Build a string and returns it.
	return w.String()
}

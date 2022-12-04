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
	// If there are no frames in the stack trace, return "No stack trace" for
	// better readability.
	if len(s.GetFrames()) == 0 {
		return "No stack trace"
	}

	// Create a new string builder and set its verbosity.
	b := toStringBuilder{Verbose: verbose}

	// Append a stack trace heading to the string builder.
	b.WriteString("Stack trace:")

	// Iterate over all frames in the stack trace and append them to the string
	// builder in a human-readable format.
	for _, f := range s.GetFrames() {
		b.appendStackFrame(f)
	}

	// Build the final string and return it.
	return b.String()
}

// toStringBuilder is a string builder that is used to build a human-readable
// string representing a stack trace.
type toStringBuilder struct {
	strings.Builder

	// Verbose indicates whether the stack trace string should include
	// more detailed information such as fully qualified function names and
	// full file paths.
	Verbose bool
}

// appendStackFrame appends a single stack frame to the string builder in a
// human-readable format.
func (b *toStringBuilder) appendStackFrame(f *stpb.StackTrace_Frame) {
	// Indent each stack frame by two spaces.
	b.WriteString("\n  ")

	// Append the function name.
	b.appendFunctionName(f)

	// Append a separator between the function name and the file name.
	b.WriteString("@")

	// Append the file name (verbose=false) or the full path
	// (verbose=true).
	b.appendFileName(f)

	// Append the line number.
	b.WriteString(":")
	b.WriteString(strconv.Itoa(int(f.GetLine())))

	// Append the program counter.
	b.appendProgramCounter(f)
}

func (b *toStringBuilder) appendFunctionName(f *stpb.StackTrace_Frame) {
	// Append the function name.
	if b.Verbose {
		// If verbose flag is enabled, append the fully qualified function
		// name.
		b.WriteString(f.GetFunction())
	} else {
		// If verbose flag is not enabled, append the short name by getting
		// the string after the last slash in the function name.
		p := strings.LastIndex(f.GetFunction(), "/")
		b.WriteString(f.GetFunction()[p+1:])
	}
}

func (b *toStringBuilder) appendFileName(f *stpb.StackTrace_Frame) {
	// Append the file name (verbose=false) or the full path
	// (verbose=true).
	if b.Verbose {
		// If verbose flag is enabled, append the full path of the file.
		b.WriteString(f.GetFile())
	} else {
		// If verbose flag is not enabled, append the base name of the
		// file.
		b.WriteString(path.Base(f.GetFile()))
	}
}

func (b *toStringBuilder) appendProgramCounter(f *stpb.StackTrace_Frame) {
	// Append the program counter if verbose=true and program counter is
	// not 0.
	if b.Verbose && f.GetProgramCounter() != 0 {
		b.WriteString("(")

		if f.GetEntry() != 0 {
			b.WriteString(fmt.Sprintf("0x%x+", f.GetEntry()))
		}

		b.WriteString(fmt.Sprintf(
			"0x%x)", f.GetProgramCounter()-f.GetEntry()))
	}
}

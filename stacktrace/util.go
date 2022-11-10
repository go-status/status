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
	if len(s.GetFrames()) == 0 {
		return "No stack trace"
	}
	var w strings.Builder
	w.WriteString("Stack trace:")
	for _, frame := range s.GetFrames() {
		w.WriteString("\n  ")
		if verbose {
			w.WriteString(frame.GetFunction())
		} else {
			p := strings.LastIndex(frame.GetFunction(), "/")
			w.WriteString(frame.GetFunction()[p+1:])
		}
		w.WriteString("@")
		if verbose {
			w.WriteString(frame.GetFile())
		} else {
			w.WriteString(path.Base(frame.GetFile()))
		}
		w.WriteString(":")
		w.WriteString(strconv.Itoa(int(frame.GetLine())))
		if verbose && frame.GetProgramCounter() != 0 {
			w.WriteString("(")
			if frame.GetEntry() != 0 {
				w.WriteString(fmt.Sprintf("0x%x+", frame.GetEntry()))
			}
			w.WriteString(fmt.Sprintf(
				"0x%x)", frame.GetProgramCounter()-frame.GetEntry()))
		}
	}
	return w.String()
}

package stacktrace

import (
	"fmt"
	"io"
	"runtime"

	stpb "github.com/go-status/status/proto"
)

// StackTrace stores a stack trace with little overhead.
type StackTrace struct {
	// A stack trace converted from a Protocol Buffers message.
	proto *stpb.StackTrace

	// A stack trace generated internally.
	// CAVEAT: The program counters are incremented by one due to the spec of
	// runtime.Callers.  Please do not depend on the implementation.
	pcs []uintptr
}

// Ensure StackTrace implements Formatter and Stringer.
var _ fmt.Formatter = StackTrace{}
var _ fmt.Stringer = StackTrace{}

// New returns the current stack trace.  Its first stack frame should identify
// the caller of New.
func New() StackTrace {
	var result StackTrace
	// NOTE: `buf` uses a constant size to avoid heap allocation.
	var buf [128]uintptr
	// NOTE: 2 represents the two stack frames: New and runtime.Callers.
	for skip := 2; ; skip += len(buf) {
		// Fetch callers and resize `pcs`.
		pcs := buf[:runtime.Callers(skip, buf[:])]
		// Append the callers to the result.
		result.pcs = append(result.pcs, pcs...)
		// If it reaches the end of callers.
		if len(pcs) < len(buf) {
			break
		}
	}

	return result
}

// NewFromProto returns a StackTrace object storing the given StackTrace proto
// message.
func NewFromProto(s *stpb.StackTrace) StackTrace {
	return StackTrace{proto: s}
}

// ToProto converts the StackTrace object consisting of program counters to
// a proto message having comprehensible stack frames (e.g., including function
// names).
func (s StackTrace) ToProto() *stpb.StackTrace {
	if s.proto != nil {
		return s.proto
	}

	result := &stpb.StackTrace{}
	for _, pc := range s.pcs {
		// NOTE: Due to the runtime.Callers spec, program counters stored in
		// StackTrace are incremented by 1.  The decrement restores a true
		// program counter.
		pc--
		frame := &stpb.StackTrace_Frame{
			File:           "unknown",
			Function:       "unknown",
			ProgramCounter: uint64(pc),
		}
		// Fill fields of the stack frame.
		if fn := runtime.FuncForPC(pc); fn != nil {
			file, line := fn.FileLine(pc)
			frame.File = file
			frame.Line = int32(line)
			frame.Function = fn.Name()
			frame.Entry = uint64(fn.Entry())
		}
		result.Frames = append(result.GetFrames(), frame)
	}
	return result
}

// Format implements fmt.Formatter.  It outputs a stack trace in a
// human-readable format using ToString. The format may change, so callers must
// not depend on it.  "%s" outputs the stack trace in a short format, and "%v"
// outputs the stack trace in a long format.
func (s StackTrace) Format(f fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = io.WriteString(f, ToString(s.ToProto(), false))
	case 'v':
		_, _ = io.WriteString(f, ToString(s.ToProto(), true))
	default:
		_, _ = io.WriteString(f, fmt.Sprintf("%%!%c(StackTrace)", verb))
	}
}

// String implements fmt.Stringer.  It returns the same format as ToString
// returns where verbose is set to false.
func (s StackTrace) String() string {
	return ToString(s.ToProto(), false)
}

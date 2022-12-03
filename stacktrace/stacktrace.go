package stacktrace

import (
	"context"
	"fmt"
	"io"
	"runtime"

	stpb "github.com/go-status/status/proto"
)

// StackTrace stores a stack trace with little overhead.
type StackTrace struct {
	// proto is a stack trace converted from a Protocol Buffers message.
	// CAVEAT: If proto is not nil, prev must be nil because proto should
	// squash stack frames.
	proto *stpb.StackTrace

	// pcs is a stack trace generated internally.
	// CAVEAT: The program counters in this field are incremented by one due
	// to the specification of runtime.Callers.  It is not recommended to
	// depend on this implementation.
	pcs []uintptr

	// prev stores the previous stack trace in the call chain.
	prev *StackTrace
}

// Ensure StackTrace implements Formatter and Stringer interfaces.
var (
	_ fmt.Formatter = &StackTrace{}
	_ fmt.Stringer  = &StackTrace{}
)

// New returns the current stack trace.  The first stack frame in the returned
// StackTrace should identify the caller of this function.
func New(ctx context.Context) StackTrace {
	result := StackTrace{prev: fromContext(ctx)}

	// NOTE: `buf` uses a constant size to avoid heap allocation.
	var buf [128]uintptr

	// NOTE: This for loop starts with a skip count of 2, to ignore the frames
	// for this function and runtime.Callers.
	for skip := 2; ; skip += len(buf) {
		// Fetch the callers and resize the pcs field in the result.
		pcs := buf[:runtime.Callers(skip, buf[:])]
		// Append the callers to the result.
		result.pcs = append(result.pcs, pcs...)

		// If the end of the callers is reached, exit the loop.
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
// a proto message with more comprehensive stack frames (e.g., including
// function names and file paths).
func (s *StackTrace) ToProto() *stpb.StackTrace {
	if s.proto != nil {
		return s.proto
	}

	result := &stpb.StackTrace{}

	s.appendTo(result)

	return result
}

// appendTo appends the stack frames from the current StackTrace object to the
// given proto message.  If the current StackTrace object has a proto message,
// the frames in the proto message are appended to the result.  Otherwise, the
// program counters in the StackTrace object are resolved to comprehensible
// stack frames and appended to the result.
func (s *StackTrace) appendTo(st *stpb.StackTrace) {
	// If the current StackTrace object has a proto message, append the frames
	// from the proto message to the result.
	if s.proto != nil {
		st.Frames = append(st.GetFrames(), s.proto.GetFrames()...)

		return
	}

	// Record the current number of frames in the result proto message.
	numFrames := len(st.GetFrames())

	for _, pc := range s.pcs {
		// NOTE: Due to the runtime.Callers spec, program counters stored in
		// StackTrace are incremented by 1.  The decrement restores a true
		// program counter.
		pc--

		// Set default values for unresolvable frames.
		frame := &stpb.StackTrace_Frame{
			File:           "unknown",
			Function:       "unknown",
			ProgramCounter: uint64(pc),
		}

		// Fill fields of the stack frame with information from the function
		// at the given program counter.
		if fn := runtime.FuncForPC(pc); fn != nil {
			file, line := fn.FileLine(pc)
			frame.File = file
			frame.Line = int32(line)
			frame.Function = fn.Name()
			frame.Entry = uint64(fn.Entry())

			// If the function name is "goEnter", discard remaining stack
			// frames because it should be treated as an entry point of
			// a go-routine.
			if fn.Name() == goEnterFuncName {
				break
			}

			// If the function name is "goExit", discard extra stack frames
			// because it is a mark to stop recording stack frames.
			if fn.Name() == goExitFuncName {
				st.Frames = st.GetFrames()[:numFrames]

				continue
			}
		}

		// Append the stack frame to the result proto message.
		st.Frames = append(st.GetFrames(), frame)
	}

	// If there is a previous stack trace, append its frames to the result as
	// well.
	if s.prev != nil {
		s.prev.appendTo(st)
	}
}

// Format implements fmt.Formatter.
// It outputs a stack trace in a human-readable format using ToString. The
// format may change, so callers must not depend on it.  "%s" outputs the
// stack trace in a short format, and "%v" outputs the stack trace in a long
// format.
func (s *StackTrace) Format(f fmt.State, verb rune) {
	switch verb {
	case 's':
		// Write the short format stack trace to the output.
		_, _ = io.WriteString(f, ToString(s.ToProto(), false))
	case 'v':
		// Write the long format stack trace to the output.
		_, _ = io.WriteString(f, ToString(s.ToProto(), true))
	default:
		// Write an error message for unsupported formats.
		_, _ = io.WriteString(f, fmt.Sprintf("%%!%c(StackTrace)", verb))
	}
}

// String implements fmt.Stringer.
// It returns the same format as ToString returns where verbose is set to
// false.
func (s *StackTrace) String() string {
	return ToString(s.ToProto(), false)
}

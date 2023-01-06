package stacktrace

import (
	"context"
	"strings"
	"sync"
	"testing"
)

// TestGo tests the Go function.
func TestGo(t *testing.T) {
	// Allow the test to be run in parallel with other tests.
	t.Parallel()

	var st *StackTrace

	var wg sync.WaitGroup

	wg.Add(1)

	// Run the Go function with a function that creates a new stack trace and
	// sets it to the st variable.
	Go(context.Background(), func(ctx context.Context) {
		tmpSt := New(ctx, 0)
		st = &tmpSt
		wg.Done()
	})

	wg.Wait()

	frames := st.ToProto().GetFrames()

	// Check each frame against the expected function suffix.
	for i, expected := range []string{
		"/stacktrace.TestGo.func1",
		"/stacktrace.Go.func1",
		"/stacktrace.Go",
		"/stacktrace.TestGo",
	} {
		// If the frame does not have the expected function suffix, fail the
		// test.
		if !strings.HasSuffix(
			frames[i].GetFunction(), expected) {
			t.Fatalf(
				"Stack frame[%d] has an unexpected function: "+
					"expected=*%s, actual=%s",
				i, expected, frames[i].GetFunction())
		}
	}
}

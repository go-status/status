package stacktrace

import (
	"context"
	"strings"
	"sync"
	"testing"
)

func TestGo(t *testing.T) {
	t.Parallel()

	var st *StackTrace

	var wg sync.WaitGroup

	wg.Add(1)

	Go(context.Background(), func(ctx context.Context) {
		tmpSt := New(ctx)
		st = &tmpSt
		wg.Done()
	})

	wg.Wait()

	frames := st.ToProto().GetFrames()

	for i, expected := range []string{
		"/stacktrace.TestGo.func1",
		"/stacktrace.Go.func1",
		"/stacktrace.Go",
		"/stacktrace.TestGo",
	} {
		if !strings.HasSuffix(
			frames[i].GetFunction(), expected) {
			t.Fatalf(
				"Stack frame[%d] has an unexpected function: "+
					"expected=*%s, actual=%s",
				i, expected, frames[i].GetFunction())
		}
	}
}

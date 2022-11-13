package stacktrace

import (
	"context"
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
	t.Errorf("%+v", st)
}

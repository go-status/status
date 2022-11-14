package stacktrace

import (
	"context"
	"runtime"
	"strings"
)

// Go spawns a go routine calling f with contextual stack tracing.  While a
// regular go routine discards stack frames of its caller, this saves the
// caller's stack frames and provides them to a spawned go-routine via Context.
func Go(ctx context.Context, f func(context.Context)) {
	// Record the current stack trace.
	st := goExit(ctx)

	// Call the give function inside a go-routine.
	go goEnter(func() { f(contextWith(ctx, st)) })
}

// goEnterFuncName is a fully-qualified function name of goEnter.
//
//nolint:gochecknoglobals
var goEnterFuncName = func() string {
	var name string

	// Acquire a stack frame of goEnter and store the fully-qualified function
	// name of goEnter to name.
	goEnter(func() {
		pc, _, _, _ := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if details != nil {
			name = details.Name()
		}
	})

	return name
}()

// goExitFuncName is a fully-qualified function name of goExit.
//
//nolint:gochecknoglobals
var goExitFuncName = strings.ReplaceAll(
	goEnterFuncName, "goEnter", "goExit")

// goEnter is a wrapper function appeared in a stack trace.  This should be
// used for marking an entry point of a go-routine.  Stack frames before this
// function should be hidden.
//
//go:noinline
func goEnter(f func()) {
	f()
}

// goExit is a wrapper function appeared in a stack trace.  This should be
// used for marking an end of a go-routine.  Stack frames after this function
// should be hidden.
//
//go:noinline
func goExit(ctx context.Context) *StackTrace {
	st := New(ctx)

	return &st
}

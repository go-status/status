package stacktrace

import (
	"context"
	"runtime"
	"strings"
)

// Go spawns a go routine calling f with contextual stack tracing.  Unlike a
// regular go routine, this saves the caller's stack frames and provides them
// to a spawned go-routine via Context.
func Go(ctx context.Context, f func(context.Context)) {
	// Record the current stack trace.
	st := goExit(ctx)

	// Call the given function inside a go-routine.
	go goEnter(func() { f(contextWith(ctx, st)) })
}

// goEnterFuncName is the fully-qualified function name of goEnter.
//
//nolint:gochecknoglobals
var goEnterFuncName = func() string {
	var name string

	// Acquire a stack frame of goEnter and store the fully-qualified function
	// name of goEnter in the name variable.
	goEnter(func() {
		pc, _, _, _ := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if details != nil {
			name = details.Name()
		}
	})

	return name
}()

// goExitFuncName is the fully-qualified function name of goExit.
//
//nolint:gochecknoglobals
var goExitFuncName = strings.ReplaceAll(
	goEnterFuncName, "goEnter", "goExit")

// goEnter is a wrapper function that appears in a stack trace.  It is used to
// mark the entry point of a go-routine, and to hide stack frames before this
// function.
//
//go:noinline
func goEnter(f func()) {
	f()
}

// goExit is a wrapper function that appears in a stack trace.  It is used to
// mark the end of a go-routine, and to hide stack frames after this function.
//
//go:noinline
func goExit(ctx context.Context) *StackTrace {
	// Create a new stack trace with the given context.
	st := New(ctx, 0)

	// Return the new stack trace.
	return &st
}

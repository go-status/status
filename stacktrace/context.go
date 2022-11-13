package stacktrace

import "context"

// contextStackTraceKeyType is a type used for a context key.
type contextStackTraceKeyType int

// contextStackTraceKey is a context key for ContextStackTrace.
const contextStackTraceKey contextStackTraceKeyType = 0

// fromContext extracts a StackTrace object from the given context.  This
// returns nil if it fails to find a StackTrace object.
func fromContext(ctx context.Context) *StackTrace {
	if ctx == nil {
		return nil
	}

	if v := ctx.Value(contextStackTraceKey); v != nil {
		if st, ok := v.(*StackTrace); ok {
			return st
		}
	}

	return nil
}

// contextWith returns a new context with the given StackTrace object.
func contextWith(ctx context.Context, st *StackTrace) context.Context {
	return context.WithValue(ctx, contextStackTraceKey, st)
}

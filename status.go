// Package status enables all errors to have stack traces.  This is useful
// because the pkg/errors package does not enforce the inclusion of stack
// traces.
package status

// Status is an error interface that contains an error code, error message,
// and stack trace.   It can also have an error chain.  In addition, it
// implements the built-in error interface so that it can be down-casted to
// an error type.  A nil value of Status represents a non-error (i.e. OK in the
// gRPC context).
// NOTE: Status does not implement methods other than Error() because calling
// a method of a nil Status value will cause a panic.
type Status interface {
	// Implements the built-in error interface to allow casting to error.
	error

	// statusImpl returns the inner statusImpl object and forces any Status to
	// embed statusImpl.
	statusImpl() *statusImpl
}

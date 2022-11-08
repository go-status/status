// Package status enables all errors to have stack traces while pkg/errors
// cannot enforce it.
package status

// Status is an error interface having an error code, an error message and
// a stack trace.  It can also have an error chain.  As well as the error
// interface, nil represents non-error (i.e., OK in the gRPC context).
// NOTE: Status does not implement methods other than Error() because Status
// can be nil and calling a method of nil panics.
type Status interface {
	// Implements the builtin `error` interface so that Status can be
	// down-casted to error.
	error

	// statusImpl returns the inner statusImpl object.  Besides, it forces
	// any Status to embed statusImpl.
	statusImpl() *statusImpl
}

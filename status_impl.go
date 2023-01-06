package status

import (
	"context"

	"github.com/go-status/status/stacktrace"
	"google.golang.org/genproto/googleapis/rpc/code"
	gRPCStatus "google.golang.org/genproto/googleapis/rpc/status"
)

type statusImpl struct {
	code       code.Code
	message    string
	stackTrace stacktrace.StackTrace
	gRPCStatus *gRPCStatus.Status
	cause      *Status
}

// Ensure statusImpl implements Status interface.
var _ Status = &statusImpl{}

func newStatusImpl(ctx context.Context, code code.Code, message string, skip int) *statusImpl {
	return &statusImpl{
		code:       code,
		message:    message,
		stackTrace: stacktrace.New(ctx, skip),
	}
}

func (st *statusImpl) Error() string {
	return st.message
}

func (st *statusImpl) statusImpl() *statusImpl {
	return st
}

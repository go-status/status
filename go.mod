module github.com/go-status/status

go 1.15

replace (
	github.com/go-status/status/proto => ./proto
	github.com/go-status/status/stacktrace => ./stacktrace
)

require (
	github.com/pkg/errors v0.9.1 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

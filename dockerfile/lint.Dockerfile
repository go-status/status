FROM golangci/golangci-lint
WORKDIR /go/src/github.com/go-status/status
COPY go.mod go.sum ./
RUN go mod download

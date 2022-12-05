FROM golangci/golangci-lint

# Set the working directory to the go-status/status repository.
WORKDIR /go/src/github.com/go-status/status

# Copy the go.mod and go.sum files to the working directory.  This is done for
# caching purposes to avoid downloading dependencies each time the code is
# built.
COPY go.mod go.sum ./

# Download the dependencies specified in the go.mod and go.sum files.
RUN go mod download

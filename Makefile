# Rule to run the `golangci-lint` tool in a Docker container to check the
# source code for style and other issues.
lint:
	docker build -t go-status/lint -f ./dockerfile/lint.Dockerfile .
	docker run --rm \
		--volume="$(shell pwd):/go/src/github.com/go-status/status" \
		--workdir=/go/src/github.com/go-status/status \
		go-status/lint golangci-lint run -v --timeout=5m

# Rule to format the source code.
format:
	go fmt ./...

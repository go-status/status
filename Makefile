
lint:
	docker build -t go-status/lint -f ./dockerfile/lint.Dockerfile .
	docker run --rm \
		--volume="$(shell pwd):/go/src/github.com/go-status/status" \
		--workdir=/go/src/github.com/go-status/status \
		go-status/lint golangci-lint run -v --timeout=5m

format:
	go fmt ./...

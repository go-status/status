DOCKER_ENVIRONMENT=$(shell if [ -f /.dockerenv ]; then echo docker; else echo host; fi)
DOCKER_IMAGE=github.com/go-status/status/protoc

all: stacktrace.pb.go

%.pb.go: %.proto
	$(MAKE) $*.pb.go@$(DOCKER_ENVIRONMENT)

%.pb.go@host: docker
	docker run --rm --volume="$(CURDIR):/work" --workdir=/work \
		"$(DOCKER_IMAGE)" make $*.pb.go@docker

%.pb.go@docker:
	mkdir /proto
	protoc --proto_path=. --go_out=/proto ./$*.proto
	mv /proto/github.com/go-status/status/proto/$*.pb.go /work/$*.pb.go

.PHONY: docker
docker: Dockerfile
	docker build -t "$(DOCKER_IMAGE)" .

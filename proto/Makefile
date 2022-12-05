# DOCKER_ENVIRONMENT is a variable that represents whether the current
# environment is Docker (docker) or not (host). This is determined by checking
# if the file /.dockerenv exists.
DOCKER_ENVIRONMENT=$(shell \
	if [ -f /.dockerenv ]; then echo docker; else echo host; fi)

# DOCKER_IMAGE is a variable that represents the Docker image used for building
# the protobuf files.
DOCKER_IMAGE=github.com/go-status/status/protoc

###############################################################################
# Main rules
###############################################################################

# Default rule to build compiled protobuf files.
all: stacktrace.pb.go

# Rule for building a .pb.go file from .proto file.
%.pb.go: %.proto
	$(MAKE) $*.pb.go@$(DOCKER_ENVIRONMENT)

# Rule for building a .pb.go file on the host machine.
%.pb.go@host: docker
	docker run --rm --volume="$(CURDIR):/work" --workdir=/work \
		"$(DOCKER_IMAGE)" make $*.pb.go@docker

# Rule for building a .pb.go file inside a Docker container.
%.pb.go@docker:
	mkdir /proto
	protoc --proto_path=. --go_out=/proto ./$*.proto
	mv /proto/github.com/go-status/status/proto/$*.pb.go /work/$*.pb.go

# Rule for building the Docker container.
.PHONY: docker
docker: Dockerfile
	docker build -t "$(DOCKER_IMAGE)" .

FROM golang:1.17-buster

# Install required dependencies.
RUN apt-get update && apt-get install -y curl unzip git make

# Download and extract the protoc binary.
RUN PROTOC_ZIP=protoc-3.17.3-linux-x86_64.zip && \
    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/$PROTOC_ZIP && \
    unzip -o $PROTOC_ZIP -d /usr/local bin/protoc && \
    chmod +x /usr/local/bin/protoc && \
    unzip -o $PROTOC_ZIP -d /usr/local 'include/*' && \
    rm -f $PROTOC_ZIP

# Make protoc available in the PATH.
ENV PATH="/usr/local/bin:$PATH"

# Install the protoc-gen-go plugin, which allows protoc to generate Go code.
RUN go install github.com/golang/protobuf/protoc-gen-go@v1.5.2

FROM golang:1.22.2

# IMPORTANT:
# The CWD for this Dockerfile is the root of the repository, not ./spec directory

WORKDIR /directmq

RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    g++ \
    gdb \
    catch2 \
    protobuf-compiler \
    libwebsockets-dev

RUN apt-get remove -y cmake && \
    wget https://github.com/Kitware/CMake/releases/download/v3.26.0/cmake-3.26.0-linux-x86_64.sh && \
    chmod +x cmake-3.26.0-linux-x86_64.sh && \
    ./cmake-3.26.0-linux-x86_64.sh --prefix=/usr/local --skip-license && \
    rm cmake-3.26.0-linux-x86_64.sh

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
RUN go install github.com/onsi/ginkgo/v2/ginkgo@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /directmq/sdk/go
COPY sdk/go/go.mod sdk/go/go.sum ./
RUN go mod download

WORKDIR /directmq/spec
COPY spec/go.mod spec/go.sum ./
RUN go mod download

WORKDIR /directmq/spec/agents/go
COPY spec/agents/go/go.mod spec/agents/go/go.sum ./
RUN go mod download

WORKDIR /directmq/spec

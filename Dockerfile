FROM ubuntu:24.04

WORKDIR /workspace

RUN apt update && apt install -y \
    build-essential \
    curl \
    gcc \
    g++ \
    gdb \
    catch2 \
    cmake \
    golang-go \
    python3 \
    python3-setuptools \
    protobuf-compiler \
    libwebsockets-dev \
    clang-tools \
    clang-format \
    clang-tidy \
    valgrind

RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
RUN mv ./bin/task /usr/local/bin/task

# Clean /workspace directory in preparation to run on-create.sh script
RUN rm -rf /workspace/*

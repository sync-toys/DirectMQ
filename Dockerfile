FROM ubuntu:24.04

WORKDIR /workspace

RUN apt update && apt install -y \
    curl \
    gcc \
    gdb \
    catch2 \
    cmake \
    clang-tools \
    python3 \
    python3-setuptools \
    protobuf-compiler

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
RUN mv ./bin/task /usr/local/bin/task

# Clean /workspace directory in preparation to run on-create.sh script
RUN rm -rf /workspace/*

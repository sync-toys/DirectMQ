FROM golang:1.22.2

# IMPORTANT:
# The CWD for this Dockerfile is the root of the repository, not ./spec directory
# so all paths are relative to the root of the repository,
# also it means that docker build context should be the root of the repository

WORKDIR /directmq

RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    g++ \
    gdbserver \
    catch2 \
    protobuf-compiler \
    python3 \
    python3-setuptools \
    python3-pip


RUN apt-get remove -y cmake && \
    wget https://github.com/Kitware/CMake/releases/download/v3.26.0/cmake-3.26.0-linux-x86_64.sh && \
    chmod +x cmake-3.26.0-linux-x86_64.sh && \
    ./cmake-3.26.0-linux-x86_64.sh --prefix=/usr/local --skip-license && \
    rm cmake-3.26.0-linux-x86_64.sh

RUN pip3 install conan --break-system-packages
RUN conan profile detect

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
RUN mv ./bin/task /usr/local/bin/task

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

# Heavy C++ dependencies fetch and pre-build runs first
WORKDIR /directmq/spec/agents/cpp
COPY spec/agents/cpp/conanfile.py spec/agents/cpp/taskfile.yml ./
RUN task conan:install

WORKDIR /directmq/sdk/cpp
COPY sdk/cpp/conanfile.py sdk/cpp/taskfile.yml ./
RUN task conan:install BUILD_TYPE=Debug
RUN task conan:install BUILD_TYPE=Release

# Now we are installing the cmake dependencies, but those require access to source
# code, so we are executing it at the end as it is the most likely to change.
# CMake will fetch its dependencies at this stage.
WORKDIR /directmq/sdk/cpp
COPY sdk/cpp/CMakeLists.txt sdk/cpp/CMakeUserPresets.json ./
COPY sdk/cpp/include/ ./include
COPY sdk/cpp/src/ ./src
COPY sdk/cpp/lib/ ./lib
RUN task cmake:generate BUILD_TYPE=Debug DMQ_BUILD_TESTS=OFF

WORKDIR /directmq/spec/agents/cpp
COPY spec/agents/cpp/*.hpp spec/agents/cpp/*.cpp spec/agents/cpp/CMakeLists.txt spec/agents/cpp/CMakeUserPresets.json ./
RUN task cmake:generate

WORKDIR /directmq/spec

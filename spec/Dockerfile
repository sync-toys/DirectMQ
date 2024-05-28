FROM golang:1.22.2

# IMPORTANT:
# The CWD for this Dockerfile is the root of the repository, not ./spec directory

WORKDIR /directmq

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
RUN go install github.com/onsi/ginkgo/v2/ginkgo@2.17.2
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

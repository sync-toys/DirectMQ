module github.com/sync-toys/DirectMQ/spec/agents/go

go 1.22.2

replace github.com/sync-toys/DirectMQ/sdk/go => ../../../sdk/go

replace github.com/sync-toys/DirectMQ/spec => ../../../spec

require (
	github.com/sync-toys/DirectMQ/sdk/go v0.0.0
	github.com/sync-toys/DirectMQ/spec v0.0.0
)

require (
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/google/pprof v0.0.0-20240424215950-a892ee059fd6 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/onsi/ginkgo/v2 v2.17.2 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/tools v0.20.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

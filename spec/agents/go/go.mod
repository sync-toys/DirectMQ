module github.com/sync-toys/DirectMQ/spec/agents/go

go 1.22.2

replace github.com/sync-toys/DirectMQ/sdk/go => ../../../sdk/go

replace github.com/sync-toys/DirectMQ/spec => ../../../spec

require (
	github.com/sync-toys/DirectMQ/sdk/go v0.0.0
	github.com/sync-toys/DirectMQ/spec v0.0.0
)

require (
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	golang.org/x/net v0.28.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)

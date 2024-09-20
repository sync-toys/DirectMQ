module github.com/sync-toys/DirectMQ/spec

go 1.22

toolchain go1.22.2

replace github.com/sync-toys/DirectMQ/sdk/go => ../sdk/go

require (
	github.com/fatih/color v1.16.0
	github.com/gorilla/websocket v1.5.1
	github.com/onsi/ginkgo/v2 v2.20.2
	github.com/onsi/gomega v1.34.1
	github.com/phayes/freeport v0.0.0-20220201140144-74d24b5ae9f5
	github.com/sync-toys/DirectMQ/sdk/go v0.0.0
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20240827171923-fa2c70bbbfe5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/tools v0.24.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

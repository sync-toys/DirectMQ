package dmqspecagents

import (
	"fmt"
	"strconv"

	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
)

const NO_DEBUGGING = -1
const DEFAULT_GO_SDK_DEBUG_PORT = 2345

func GolangAgent(nodeID string, debugPort int) dmqspecagent.UniversalSpawn {
	executablePath := "/directmq/spec/agents/go/bin/go-agent"

	if debugPort == NO_DEBUGGING {
		return dmqspecagent.UniversalSpawn{
			NodeID:         nodeID,
			ExecutablePath: executablePath,
			Arguments:      []string{},
		}
	}

	fmt.Println("debugPort: ", debugPort)
	return dmqspecagent.UniversalSpawn{
		NodeID:         nodeID,
		ExecutablePath: "dlv",
		Arguments:      []string{"exec", executablePath, "--headless", "--listen=0.0.0.0:" + strconv.Itoa(debugPort), "--api-version=2"},
	}
}

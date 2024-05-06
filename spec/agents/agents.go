package dmqspecagents

import (
	"strconv"

	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
)

const NO_DEBUGGING = -1

func GolangAgent(nodeID string, debugPort int) dmqspecagent.UniversalSpawn {
	executablePath := "/directmq/spec/agents/go/bin/go-agent"

	if debugPort == NO_DEBUGGING {
		return dmqspecagent.UniversalSpawn{
			NodeID:         nodeID,
			ExecutablePath: executablePath,
			Arguments:      []string{},
		}
	}

	return dmqspecagent.UniversalSpawn{
		NodeID:         nodeID,
		ExecutablePath: "dlv",
		Arguments:      []string{"exec", executablePath, "--headless", "--listen=:" + strconv.Itoa(debugPort), "--api-version=2"},
	}
}

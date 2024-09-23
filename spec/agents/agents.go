package dmqspecagents

import (
	"strconv"

	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
)

const NO_DEBUGGING = -1
const DEFAULT_GO_SDK_DEBUG_PORT = 2345
const DEFAULT_CPP_SDK_DEBUG_PORT = 3456

func GetDefaultDebugPort(agentType string) int {
	switch agentType {
	case "go":
		return DEFAULT_GO_SDK_DEBUG_PORT
	case "cpp":
		return DEFAULT_CPP_SDK_DEBUG_PORT
	default:
		panic("Unknown agent type, cannot determine default debug port")
	}
}

func GolangAgent(nodeID string, debugPort int) dmqspecagent.UniversalSpawn {
	executablePath := "/directmq/spec/agents/go/bin/go-agent"

	if debugPort == NO_DEBUGGING {
		return dmqspecagent.UniversalSpawn{
			AgentType:      "go",
			NodeID:         nodeID,
			ExecutablePath: executablePath,
			Arguments:      []string{},
		}
	}

	return dmqspecagent.UniversalSpawn{
		AgentType:      "go",
		NodeID:         nodeID,
		ExecutablePath: "dlv",
		Arguments:      []string{"exec", executablePath, "--headless", "--listen=0.0.0.0:" + strconv.Itoa(debugPort), "--api-version=2"},
	}
}

func CppAgent(nodeID string, debugPort int) dmqspecagent.UniversalSpawn {
	executablePath := "/directmq/spec/agents/cpp/build/directmq-sdk-agent"

	if debugPort == NO_DEBUGGING {
		return dmqspecagent.UniversalSpawn{
			AgentType:      "cpp",
			NodeID:         nodeID,
			ExecutablePath: executablePath,
			Arguments:      []string{},
		}
	}

	return dmqspecagent.UniversalSpawn{
		AgentType:      "cpp",
		NodeID:         nodeID,
		ExecutablePath: "gdbserver",
		Arguments:      []string{"0.0.0.0:" + strconv.Itoa(debugPort), executablePath},
	}
}

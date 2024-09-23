package dmqspecagents

import (
	"fmt"
	"math"
	"strings"

	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
)

type UniversalSpawner func(nodeID string, debugPort int) dmqspecagent.UniversalSpawn

type combinations struct {
	spawners []UniversalSpawner
}

func newCombinations(spawners []UniversalSpawner) *combinations {
	return &combinations{spawners}
}

var Combinations = newCombinations([]UniversalSpawner{GolangAgent, CppAgent})

type combination struct {
	nodeId  string
	spawner UniversalSpawner
}

func (p *combinations) generateSpawnerCombinations(nodeIds []string) [][]combination {
	firstListLen := len(nodeIds)
	secondListLen := len(p.spawners)

	combinationCount := int(math.Pow(float64(secondListLen), float64(firstListLen)))

	combinations := make([][]combination, combinationCount)
	combinationRow := make([]combination, firstListLen)

	for i := 0; i < combinationCount; i++ {
		for j := 0; j < firstListLen; j++ {
			index := (i / int(math.Pow(float64(secondListLen), float64(j)))) % secondListLen
			combinationRow[j] = combination{nodeIds[j], p.spawners[index]}
		}
		combinations[i] = combinationRow
	}

	return combinations
}

func (p *combinations) generateNamedSpawnerCombinations(nodeIds ...string) []map[string]UniversalSpawner {
	combinations := p.generateSpawnerCombinations(nodeIds)
	result := make([]map[string]UniversalSpawner, len(combinations))

	for i, combination := range combinations {
		result[i] = make(map[string]UniversalSpawner)
		for _, c := range combination {
			result[i][c.nodeId] = c.spawner
		}
	}

	return result
}

func (p *combinations) Generate(nodeIds ...string) []map[string]dmqspecagent.UniversalSpawn {
	namedCombinations := p.generateNamedSpawnerCombinations(nodeIds...)
	result := make([]map[string]dmqspecagent.UniversalSpawn, len(namedCombinations))

	for i, combination := range namedCombinations {
		result[i] = make(map[string]dmqspecagent.UniversalSpawn)
		for nodeId, spawner := range combination {
			result[i][nodeId] = spawner(nodeId, NO_DEBUGGING)
		}
	}

	return result
}

func (p *combinations) Describe(combinations []map[string]dmqspecagent.UniversalSpawn) []string {
	result := make([]string, len(combinations))

	for i, combination := range combinations {
		test := make([]string, 0)
		for nodeId, spawn := range combination {
			test = append(test, fmt.Sprintf("%s: %s", nodeId, spawn.AgentType))
		}

		result[i] = strings.Join(test, ", ")
	}

	return result
}

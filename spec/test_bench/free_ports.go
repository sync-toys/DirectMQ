package testbench

import (
	"sync"

	"github.com/phayes/freeport"
)

var usedPorts = make([]int, 0)
var usedPortsMutex = sync.Mutex{}

func getFreeTestingPort() int {
	usedPortsMutex.Lock()
	defer usedPortsMutex.Unlock()

	port, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}

	usedPorts = append(usedPorts, port)
	return port
}

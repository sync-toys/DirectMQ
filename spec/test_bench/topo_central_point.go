package testbench

import (
	"fmt"
	"strconv"
	"time"

	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
)

type CentralPointTopoTestBenchConfig struct {
	TopSpawn     dmqspecagent.UniversalSpawn
	LeftSpawn    dmqspecagent.UniversalSpawn
	CentralSpawn dmqspecagent.UniversalSpawn
	RightSpawn   dmqspecagent.UniversalSpawn

	TopTTL            int32
	TopMaxMessageSize uint64

	LeftTTL            int32
	LeftMaxMessageSize uint64

	CentralTTL            int32
	CentralMaxMessageSize uint64

	RightTTL            int32
	RightMaxMessageSize uint64

	LogTopToCentralCommunication bool
	LogCentralToTopCommunication bool

	LogLeftToCentralCommunication bool
	LogCentralToLeftCommunication bool

	LogCentralToRightCommunication bool
	LogRightToCentralCommunication bool

	LogTopLogs     bool
	LogLeftLogs    bool
	LogCentralLogs bool
	LogRightLogs   bool

	DisableAllLogs bool
}

type CentralPointTopoTestBenchTestFrameworkIntegration struct {
	CentralPointTopoTestBenchConfig

	AgentLogger          AgentLogger
	RecordingsComparator dmqspecagent.RecordingsComparator
	RecordingsLogger     RecordingsLogger
	BenchLogger          BenchLogger
}

type CentralPointTopoTestBench struct {
	config CentralPointTopoTestBenchTestFrameworkIntegration

	Top     dmqspecagent.Agent
	Left    dmqspecagent.Agent
	Central dmqspecagent.Agent
	Right   dmqspecagent.Agent

	topCentralSpy   *NodesCommunicationSpy
	leftCentralSpy  *NodesCommunicationSpy
	rightCentralSpy *NodesCommunicationSpy

	topReady     chan struct{}
	topConnected chan struct{}

	leftReady     chan struct{}
	leftConnected chan struct{}

	centralReady     chan struct{}
	centralConnected chan string

	rightReady     chan struct{}
	rightConnected chan struct{}

	alreadyStarted bool
	alreadyStopped bool
}

func NewCentralPointTopoTestBench(config CentralPointTopoTestBenchTestFrameworkIntegration) *CentralPointTopoTestBench {
	validateCentralPointTopoTestBenchConfig(config)

	return &CentralPointTopoTestBench{
		config: config,

		Left:    nil,
		Central: nil,
		Right:   nil,

		topCentralSpy:   nil,
		leftCentralSpy:  nil,
		rightCentralSpy: nil,

		topReady:     make(chan struct{}),
		topConnected: make(chan struct{}),

		leftReady:     make(chan struct{}),
		leftConnected: make(chan struct{}),

		centralReady:     make(chan struct{}),
		centralConnected: make(chan string),

		rightReady:     make(chan struct{}),
		rightConnected: make(chan struct{}),

		alreadyStarted: false,
		alreadyStopped: false,
	}
}

func NewGinkgoCentralPointTopoTestBench(config CentralPointTopoTestBenchConfig) *CentralPointTopoTestBench {
	return NewCentralPointTopoTestBench(CentralPointTopoTestBenchTestFrameworkIntegration{
		CentralPointTopoTestBenchConfig: config,

		BenchLogger:          GinkgoBenchLogger,
		AgentLogger:          GinkgoAgentLogger,
		RecordingsComparator: GomegaRecordingsComparator,
		RecordingsLogger:     GinkgoRecordingsLogger,
	})
}

func (bench *CentralPointTopoTestBench) Start() {
	defer bench.forceCleanUpTestBenchOnPanic()

	if bench.alreadyStarted {
		panic("test bench already started")
	}

	bench.benchLog("starting test bench")
	bench.alreadyStarted = true

	bench.configure()

	bench.benchLog("starting spies")
	bench.startSpies()

	bench.benchLog("starting agents")
	bench.startAgents()

	bench.benchLog("connecting agents")
	bench.connectAgents()

	bench.benchLog("bench ready")
}

func (bench *CentralPointTopoTestBench) Stop(reason string) {
	defer bench.forceCleanUpTestBenchOnPanic()

	if bench.alreadyStopped {
		bench.benchLog("test bench already stopped, ignoring stop request")
		return
	}

	bench.benchLog("stopping test bench")
	bench.alreadyStopped = true

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	done := make(chan struct{})
	go func() {
		bench.stopAgents(reason)
		bench.stopSpies()
		close(done)
	}()

	select {
	case <-timer.C:
		panic("test bench stop timeout")
	case <-done:
		bench.benchLog("test bench stopped")
	}
}

func (bench *CentralPointTopoTestBench) StopWithoutTimeout(reason string) {
	defer bench.forceCleanUpTestBenchOnPanic()

	if bench.alreadyStopped {
		bench.benchLog("test bench already stopped, ignoring stop request")
		return
	}

	bench.benchLog("stopping test bench")
	bench.alreadyStopped = true

	bench.benchLog("stopping agents")
	bench.stopAgents(reason)

	bench.benchLog("stopping spies")
	bench.stopSpies()

	bench.benchLog("test bench stopped")
}

func (bench *CentralPointTopoTestBench) startSpies() {
	bench.benchLog("starting top-central spy")
	go bench.topCentralSpy.Start()

	bench.benchLog("starting left-central spy")
	go bench.leftCentralSpy.Start()

	bench.benchLog("starting central-right spy")
	go bench.rightCentralSpy.Start()
}

func (bench *CentralPointTopoTestBench) startAgents() {
	bench.benchLog("starting top agent")
	bench.Top.Run(
		dmqspecagent.UniversalSpawn{
			NodeID:         bench.config.TopSpawn.NodeID,
			ExecutablePath: bench.config.TopSpawn.ExecutablePath,
			Arguments:      bench.config.TopSpawn.Arguments,
		},
		dmqspecagent.SetupCommand{
			NodeID:         bench.config.TopSpawn.NodeID,
			TTL:            bench.config.TopTTL,
			MaxMessageSize: bench.config.TopMaxMessageSize,
		},
	)

	bench.benchLog("waiting for top agent to be ready")
	<-bench.topReady

	bench.benchLog("starting left agent")
	bench.Left.Run(
		dmqspecagent.UniversalSpawn{
			NodeID:         bench.config.LeftSpawn.NodeID,
			ExecutablePath: bench.config.LeftSpawn.ExecutablePath,
			Arguments:      bench.config.LeftSpawn.Arguments,
		},
		dmqspecagent.SetupCommand{
			NodeID:         bench.config.LeftSpawn.NodeID,
			TTL:            bench.config.LeftTTL,
			MaxMessageSize: bench.config.LeftMaxMessageSize,
		},
	)

	bench.benchLog("waiting for left agent to be ready")
	<-bench.leftReady

	bench.benchLog("starting right agent")
	bench.Right.Run(
		dmqspecagent.UniversalSpawn{
			NodeID:         bench.config.RightSpawn.NodeID,
			ExecutablePath: bench.config.RightSpawn.ExecutablePath,
			Arguments:      bench.config.RightSpawn.Arguments,
		},
		dmqspecagent.SetupCommand{
			NodeID:         bench.config.RightSpawn.NodeID,
			TTL:            bench.config.RightTTL,
			MaxMessageSize: bench.config.RightMaxMessageSize,
		},
	)

	bench.benchLog("waiting for right agent to be ready")
	<-bench.rightReady

	bench.benchLog("starting central agent")
	bench.Central.Run(
		dmqspecagent.UniversalSpawn{
			NodeID:         bench.config.CentralSpawn.NodeID,
			ExecutablePath: bench.config.CentralSpawn.ExecutablePath,
			Arguments:      bench.config.CentralSpawn.Arguments,
		},
		dmqspecagent.SetupCommand{
			NodeID:         bench.config.CentralSpawn.NodeID,
			TTL:            bench.config.CentralTTL,
			MaxMessageSize: bench.config.CentralMaxMessageSize,
		},
	)

	bench.benchLog("waiting for central agent to be ready")
	<-bench.centralReady
}

func (bench *CentralPointTopoTestBench) connectAgents() {
	bench.benchLog("listening on central agent")
	bench.Central.Listen(dmqspecagent.ListenCommand{
		Address: bench.leftCentralSpy.ToURL.String(),
	})

	bench.benchLog("connecting top agent to central agent")
	bench.Top.Connect(dmqspecagent.ConnectCommand{
		Address: bench.topCentralSpy.FromURL.String(),
	})

	bench.benchLog("connecting left agent to central agent")
	bench.Left.Connect(dmqspecagent.ConnectCommand{
		Address: bench.leftCentralSpy.FromURL.String(),
	})

	bench.benchLog("connecting right agent to central agent")
	bench.Right.Connect(dmqspecagent.ConnectCommand{
		Address: bench.rightCentralSpy.FromURL.String(),
	})

	bench.benchLog("waiting for top agent to connect")
	<-bench.topConnected

	bench.benchLog("waiting for central agent to receive connection from top agent")
	<-bench.centralConnected

	bench.benchLog("waiting for left agent to connect")
	<-bench.leftConnected

	bench.benchLog("waiting for central agent to receive connection from left agent")
	<-bench.centralConnected

	bench.benchLog("waiting for right agent to connect")
	<-bench.rightConnected

	bench.benchLog("waiting for central agent to receive connection from right agent")
	<-bench.centralConnected
}

func (bench *CentralPointTopoTestBench) forceCleanUpTestBenchOnPanic() {
	if r := recover(); r != nil {
		bench.benchLog("force cleaning up test bench due to panic")

		bench.killAgents()
		bench.stopSpies()
	}
}

func (bench *CentralPointTopoTestBench) stopAgents(reason string) {
	bench.benchLog("stopping top agent")
	bench.Top.Stop(dmqspecagent.StopCommand{
		Reason: reason,
	})

	bench.benchLog("stopping left agent")
	bench.Left.Stop(dmqspecagent.StopCommand{
		Reason: reason,
	})

	bench.benchLog("stopping right agent")
	bench.Right.Stop(dmqspecagent.StopCommand{
		Reason: reason,
	})

	bench.benchLog("stopping central agent")
	bench.Central.Stop(dmqspecagent.StopCommand{
		Reason: reason,
	})
}

func (bench *CentralPointTopoTestBench) stopSpies() {
	bench.benchLog("stopping top-central spy")
	bench.topCentralSpy.Stop()

	bench.benchLog("stopping left-central spy")
	bench.leftCentralSpy.Stop()

	bench.benchLog("stopping central-right spy")
	bench.rightCentralSpy.Stop()
}

func (bench *CentralPointTopoTestBench) killAgents() {
	bench.benchLog("force killing top agent")
	bench.Top.Kill()

	bench.benchLog("force killing left agent")
	bench.Left.Kill()

	bench.benchLog("force killing right agent")
	bench.Right.Kill()

	bench.benchLog("force killing central agent")
	bench.Central.Kill()
}

func (bench *CentralPointTopoTestBench) configure() {
	bench.configureTop()
	bench.configureLeft()
	bench.configureCentral()
	bench.configureRight()

	bench.configureTopCentralSpy()
	bench.configureLeftCentralSpy()
	bench.configureRightCentralSpy()
}

func (bench *CentralPointTopoTestBench) configureTop() {
	bench.benchLog("configuring top agent")

	bench.Top = dmqspecagent.NewUniversalAgent(bench.config.TopSpawn.NodeID)

	bench.Top.OnReady(func(_ dmqspecagent.ReadyNotification) {
		bench.topReady <- struct{}{}
	})

	bench.Top.OnConnectionEstablished(func(_ dmqspecagent.ConnectionEstablishedNotification) {
		bench.topConnected <- struct{}{}
	})

	bench.Top.OnLogEntry(func(log string) {
		if !bench.config.DisableAllLogs && bench.config.LogTopLogs {
			bench.config.AgentLogger(log, bench.Top.GetNodeID())
		}
	})

	bench.Top.OnFatal(func(fatal dmqspecagent.FatalNotification) {
		bench.config.AgentLogger("FATAL: "+fatal.Err, bench.Top.GetNodeID())
	})
}

func (bench *CentralPointTopoTestBench) configureLeft() {
	bench.benchLog("configuring left agent")

	bench.Left = dmqspecagent.NewUniversalAgent(bench.config.LeftSpawn.NodeID)

	bench.Left.OnReady(func(_ dmqspecagent.ReadyNotification) {
		bench.leftReady <- struct{}{}
	})

	bench.Left.OnConnectionEstablished(func(_ dmqspecagent.ConnectionEstablishedNotification) {
		bench.leftConnected <- struct{}{}
	})

	bench.Left.OnLogEntry(func(log string) {
		if !bench.config.DisableAllLogs && bench.config.LogLeftLogs {
			bench.config.AgentLogger(log, bench.Left.GetNodeID())
		}
	})

	bench.Left.OnFatal(func(fatal dmqspecagent.FatalNotification) {
		bench.config.AgentLogger("FATAL: "+fatal.Err, bench.Left.GetNodeID())
	})
}

func (bench *CentralPointTopoTestBench) configureCentral() {
	bench.benchLog("configuring central agent")

	bench.Central = dmqspecagent.NewUniversalAgent(bench.config.CentralSpawn.NodeID)

	bench.Central.OnReady(func(_ dmqspecagent.ReadyNotification) {
		bench.centralReady <- struct{}{}
	})

	bench.Central.OnConnectionEstablished(func(connected dmqspecagent.ConnectionEstablishedNotification) {
		bench.centralConnected <- connected.BridgedNodeId
	})

	bench.Central.OnLogEntry(func(log string) {
		if !bench.config.DisableAllLogs && bench.config.LogCentralLogs {
			bench.config.AgentLogger(log, bench.Central.GetNodeID())
		}
	})

	bench.Central.OnFatal(func(fatal dmqspecagent.FatalNotification) {
		bench.config.AgentLogger("FATAL: "+fatal.Err, bench.Central.GetNodeID())
	})
}

func (bench *CentralPointTopoTestBench) configureRight() {
	bench.benchLog("configuring right agent")

	bench.Right = dmqspecagent.NewUniversalAgent(bench.config.RightSpawn.NodeID)

	bench.Right.OnReady(func(_ dmqspecagent.ReadyNotification) {
		bench.rightReady <- struct{}{}
	})

	bench.Right.OnConnectionEstablished(func(_ dmqspecagent.ConnectionEstablishedNotification) {
		bench.rightConnected <- struct{}{}
	})

	bench.Right.OnLogEntry(func(log string) {
		if !bench.config.DisableAllLogs && bench.config.LogRightLogs {
			bench.config.AgentLogger(log, bench.Right.GetNodeID())
		}
	})

	bench.Right.OnFatal(func(fatal dmqspecagent.FatalNotification) {
		bench.config.AgentLogger("FATAL: "+fatal.Err, bench.Right.GetNodeID())
	})
}

func (bench *CentralPointTopoTestBench) configureTopCentralSpy() {
	bench.benchLog("configuring top-central spy")

	bench.topCentralSpy = NewNodesCommunicationSpy(bench.Central, bench.Top, NodesCommunicationSpyConfig{
		LogAToB: bench.config.LogCentralToTopCommunication && !bench.config.DisableAllLogs,
		LogBToA: bench.config.LogTopToCentralCommunication && !bench.config.DisableAllLogs,

		RecordingsComparator: bench.config.RecordingsComparator,
		RecordingsLogger:     bench.config.RecordingsLogger,
	})
}

func (bench *CentralPointTopoTestBench) getCentralNodePort() int {
	centralPort, err := strconv.Atoi(bench.topCentralSpy.ToURL.Port())
	if err != nil {
		panic("Failed to get central node port: " + err.Error())
	}

	return centralPort
}

func (bench *CentralPointTopoTestBench) configureLeftCentralSpy() {
	bench.benchLog("configuring left-central spy")

	bench.leftCentralSpy = NewNodesCommunicationSpy(bench.Central, bench.Left, NodesCommunicationSpyConfig{
		LogAToB: bench.config.LogCentralToLeftCommunication && !bench.config.DisableAllLogs,
		LogBToA: bench.config.LogLeftToCentralCommunication && !bench.config.DisableAllLogs,

		RecordingsComparator: bench.config.RecordingsComparator,
		RecordingsLogger:     bench.config.RecordingsLogger,

		ForwardToPort: bench.getCentralNodePort(),
	})
}

func (bench *CentralPointTopoTestBench) configureRightCentralSpy() {
	bench.benchLog("configuring right-central spy")

	bench.rightCentralSpy = NewNodesCommunicationSpy(bench.Central, bench.Right, NodesCommunicationSpyConfig{
		LogAToB: bench.config.LogCentralToRightCommunication && !bench.config.DisableAllLogs,
		LogBToA: bench.config.LogRightToCentralCommunication && !bench.config.DisableAllLogs,

		RecordingsComparator: bench.config.RecordingsComparator,
		RecordingsLogger:     bench.config.RecordingsLogger,

		ForwardToPort: bench.getCentralNodePort(),
	})
}

func (bench *CentralPointTopoTestBench) benchLog(message string) {
	if !bench.config.DisableAllLogs {
		bench.config.BenchLogger(message)
	}
}

func (bench *CentralPointTopoTestBench) SnapshotRecordings(recordingsPrefix string) {
	bench.topCentralSpy.SnapshotRecording(fmt.Sprintf(recordingsPrefix, "top_central"))
	bench.leftCentralSpy.SnapshotRecording(fmt.Sprintf(recordingsPrefix, "left_central"))
	bench.rightCentralSpy.SnapshotRecording(fmt.Sprintf(recordingsPrefix, "right_central"))
}

func validateCentralPointTopoTestBenchConfig(config CentralPointTopoTestBenchTestFrameworkIntegration) {
	nodesToCheck := []dmqspecagent.UniversalSpawn{
		config.TopSpawn,
		config.LeftSpawn,
		config.CentralSpawn,
		config.RightSpawn,
	}

	nodeIds := make([]string, 4)
	for _, node := range nodesToCheck {
		if node.NodeID == "" {
			panic("Node ID must be provided")
		}

		if contains(nodeIds, node.NodeID) {
			panic("Node IDs must be unique")
		}

		nodeIds = append(nodeIds, node.NodeID)

		if node.ExecutablePath == "" {
			panic("Executable path must be provided")
		}
	}

	if config.TopTTL <= 0 {
		panic("Top TTL must be greater than 0")
	}

	if config.LeftTTL <= 0 {
		panic("Left TTL must be greater than 0")
	}

	if config.CentralTTL <= 0 {
		panic("Central TTL must be greater than 0")
	}

	if config.RightTTL <= 0 {
		panic("Right TTL must be greater than 0")
	}

	if config.AgentLogger == nil {
		panic("Agent logger must be provided")
	}

	if config.RecordingsComparator == nil {
		panic("Recordings comparator must be provided")
	}

	if config.RecordingsLogger == nil {
		panic("Recordings logger must be provided")
	}
}

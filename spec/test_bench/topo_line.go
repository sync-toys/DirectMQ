package testbench

import (
	"strconv"
	"time"

	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
)

type LineTopoTestBenchConfig struct {
	LeftSpawn   dmqspecagent.UniversalSpawn
	MiddleSpawn dmqspecagent.UniversalSpawn
	RightSpawn  dmqspecagent.UniversalSpawn

	LeftTTL            int32
	LeftMaxMessageSize uint64

	MiddleTTL            int32
	MiddleMaxMessageSize uint64

	RightTTL            int32
	RightMaxMessageSize uint64

	LogLeftToMiddleCommunication bool
	LogMiddleToLeftCommunication bool

	LogMiddleToRightCommunication bool
	LogRightToMiddleCommunication bool

	LogLeftLogs   bool
	LogMiddleLogs bool
	LogRightLogs  bool

	DisableAllLogs bool
}

type LineTopoTestBenchTestFrameworkIntegration struct {
	LineTopoTestBenchConfig

	AgentLogger          AgentLogger
	RecordingsComparator dmqspecagent.RecordingsComparator
	RecordingsLogger     RecordingsLogger
	BenchLogger          BenchLogger
}

type LineTopoTestBench struct {
	config LineTopoTestBenchTestFrameworkIntegration

	Left   dmqspecagent.Agent
	Middle dmqspecagent.Agent
	Right  dmqspecagent.Agent

	leftMiddleSpy  *NodesCommunicationSpy
	rightMiddleSpy *NodesCommunicationSpy

	leftReady     chan struct{}
	leftConnected chan struct{}

	middleReady     chan struct{}
	middleConnected chan string

	rightReady     chan struct{}
	rightConnected chan struct{}

	alreadyStarted bool
	alreadyStopped bool
}

func NewLineTopoTestBench(config LineTopoTestBenchTestFrameworkIntegration) *LineTopoTestBench {
	validateLineTopoTestBenchConfig(config)

	return &LineTopoTestBench{
		config: config,

		Left:   nil,
		Middle: nil,
		Right:  nil,

		leftMiddleSpy:  nil,
		rightMiddleSpy: nil,

		leftReady:     make(chan struct{}),
		leftConnected: make(chan struct{}),

		middleReady:     make(chan struct{}),
		middleConnected: make(chan string),

		rightReady:     make(chan struct{}),
		rightConnected: make(chan struct{}),

		alreadyStarted: false,
		alreadyStopped: false,
	}
}

func NewGinkgoLineTopoTestBench(config LineTopoTestBenchConfig) *LineTopoTestBench {
	return NewLineTopoTestBench(LineTopoTestBenchTestFrameworkIntegration{
		LineTopoTestBenchConfig: config,

		BenchLogger:          GinkgoBenchLogger,
		AgentLogger:          GinkgoAgentLogger,
		RecordingsComparator: GomegaRecordingsComparator,
		RecordingsLogger:     GinkgoRecordingsLogger,
	})
}

func (bench *LineTopoTestBench) Start() {
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

func (bench *LineTopoTestBench) Stop(reason string) {
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

func (bench *LineTopoTestBench) startSpies() {
	bench.benchLog("starting left-middle spy")
	go bench.leftMiddleSpy.Start()

	bench.benchLog("starting middle-right spy")
	go bench.rightMiddleSpy.Start()
}

func (bench *LineTopoTestBench) startAgents() {
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

	bench.benchLog("starting middle agent")
	bench.Middle.Run(
		dmqspecagent.UniversalSpawn{
			NodeID:         bench.config.MiddleSpawn.NodeID,
			ExecutablePath: bench.config.MiddleSpawn.ExecutablePath,
			Arguments:      bench.config.MiddleSpawn.Arguments,
		},
		dmqspecagent.SetupCommand{
			NodeID:         bench.config.MiddleSpawn.NodeID,
			TTL:            bench.config.MiddleTTL,
			MaxMessageSize: bench.config.MiddleMaxMessageSize,
		},
	)

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

	bench.benchLog("waiting for left agent to be ready")
	<-bench.leftReady

	bench.benchLog("waiting for middle agent to be ready")
	<-bench.middleReady

	bench.benchLog("waiting for right agent to be ready")
	<-bench.rightReady
}

func (bench *LineTopoTestBench) connectAgents() {
	bench.benchLog("listening on middle agent")
	bench.Middle.Listen(dmqspecagent.ListenCommand{
		Address: bench.leftMiddleSpy.ToURL.String(),
	})

	bench.benchLog("connecting left agent to middle agent")
	bench.Left.Connect(dmqspecagent.ConnectCommand{
		Address: bench.leftMiddleSpy.FromURL.String(),
	})

	bench.benchLog("connecting right agent to middle agent")
	bench.Right.Connect(dmqspecagent.ConnectCommand{
		Address: bench.rightMiddleSpy.FromURL.String(),
	})

	bench.benchLog("waiting for left agent to connect")
	<-bench.leftConnected

	bench.benchLog("waiting for middle agent to receive connection from left agent")
	<-bench.middleConnected

	bench.benchLog("waiting for right agent to connect")
	<-bench.rightConnected

	bench.benchLog("waiting for middle agent to receive connection from right agent")
	<-bench.middleConnected
}

func (bench *LineTopoTestBench) forceCleanUpTestBenchOnPanic() {
	if r := recover(); r != nil {
		bench.benchLog("force cleaning up test bench due to panic")

		bench.killAgents()
		bench.stopSpies()
	}
}

func (bench *LineTopoTestBench) stopAgents(reason string) {
	bench.benchLog("stopping left agent")
	bench.Left.Stop(dmqspecagent.StopCommand{
		Reason: reason,
	})

	bench.benchLog("stopping middle agent")
	bench.Middle.Stop(dmqspecagent.StopCommand{
		Reason: reason,
	})

	bench.benchLog("stopping right agent")
	bench.Right.Stop(dmqspecagent.StopCommand{
		Reason: reason,
	})
}

func (bench *LineTopoTestBench) stopSpies() {
	bench.benchLog("stopping left-middle spy")
	bench.leftMiddleSpy.Stop()

	bench.benchLog("stopping middle-right spy")
	bench.rightMiddleSpy.Stop()
}

func (bench *LineTopoTestBench) killAgents() {
	bench.benchLog("force killing left agent")
	bench.Left.Kill()

	bench.benchLog("force killing middle agent")
	bench.Middle.Kill()

	bench.benchLog("force killing right agent")
	bench.Right.Kill()
}

func (bench *LineTopoTestBench) configure() {
	bench.configureLeft()
	bench.configureMiddle()
	bench.configureRight()
	bench.configureLeftMiddleSpy()
	bench.configureMiddleRightSpy()
}

func (bench *LineTopoTestBench) configureLeft() {
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

func (bench *LineTopoTestBench) configureMiddle() {
	bench.benchLog("configuring middle agent")

	bench.Middle = dmqspecagent.NewUniversalAgent(bench.config.MiddleSpawn.NodeID)

	bench.Middle.OnReady(func(_ dmqspecagent.ReadyNotification) {
		bench.middleReady <- struct{}{}
	})

	bench.Middle.OnConnectionEstablished(func(connected dmqspecagent.ConnectionEstablishedNotification) {
		bench.middleConnected <- connected.BridgedNodeId
	})

	bench.Middle.OnLogEntry(func(log string) {
		if !bench.config.DisableAllLogs && bench.config.LogMiddleLogs {
			bench.config.AgentLogger(log, bench.Middle.GetNodeID())
		}
	})

	bench.Middle.OnFatal(func(fatal dmqspecagent.FatalNotification) {
		bench.config.AgentLogger("FATAL: "+fatal.Err, bench.Middle.GetNodeID())
	})
}

func (bench *LineTopoTestBench) configureRight() {
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

func (bench *LineTopoTestBench) configureLeftMiddleSpy() {
	bench.benchLog("configuring left-middle spy")

	bench.leftMiddleSpy = NewNodesCommunicationSpy(bench.Left, bench.Middle, NodesCommunicationSpyConfig{
		LogAToB: bench.config.LogLeftToMiddleCommunication && !bench.config.DisableAllLogs,
		LogBToA: bench.config.LogMiddleToLeftCommunication && !bench.config.DisableAllLogs,

		RecordingsComparator: bench.config.RecordingsComparator,
		RecordingsLogger:     bench.config.RecordingsLogger,
	})
}

func (bench *LineTopoTestBench) configureMiddleRightSpy() {
	bench.benchLog("configuring middle-right spy")

	middleNodePort, err := strconv.Atoi(bench.leftMiddleSpy.ToURL.Port())
	if err != nil {
		panic("Failed to get middle node port: " + err.Error())
	}

	bench.rightMiddleSpy = NewNodesCommunicationSpy(bench.Middle, bench.Right, NodesCommunicationSpyConfig{
		LogAToB: bench.config.LogMiddleToRightCommunication && !bench.config.DisableAllLogs,
		LogBToA: bench.config.LogRightToMiddleCommunication && !bench.config.DisableAllLogs,

		RecordingsComparator: bench.config.RecordingsComparator,
		RecordingsLogger:     bench.config.RecordingsLogger,

		ForwardToPort: middleNodePort,
	})
}

func (bench *LineTopoTestBench) benchLog(message string) {
	if !bench.config.DisableAllLogs {
		bench.config.BenchLogger(message)
	}
}

func validateLineTopoTestBenchConfig(config LineTopoTestBenchTestFrameworkIntegration) {
	if config.LeftSpawn.NodeID == config.MiddleSpawn.NodeID {
		panic("Left and middle node IDs must be different")
	}

	if config.MiddleSpawn.NodeID == config.RightSpawn.NodeID {
		panic("Middle and right node IDs must be different")
	}

	if config.LeftSpawn.NodeID == config.RightSpawn.NodeID {
		panic("Left and right node IDs must be different")
	}

	if config.LeftTTL <= 0 {
		panic("Left TTL must be greater than 0")
	}

	if config.MiddleTTL <= 0 {
		panic("Middle TTL must be greater than 0")
	}

	if config.RightTTL <= 0 {
		panic("Right TTL must be greater than 0")
	}

	if config.LeftMaxMessageSize < 0 {
		panic("Left max message size must be greater than or equal to 0")
	}

	if config.MiddleMaxMessageSize < 0 {
		panic("Middle max message size must be greater than or equal to 0")
	}

	if config.RightMaxMessageSize < 0 {
		panic("Right max message size must be greater than or equal to 0")
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

package testbench

import (
	"time"

	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
)

type LoopTopoTestBenchConfig struct {
	LeftSpawn  dmqspecagent.UniversalSpawn
	TopSpawn   dmqspecagent.UniversalSpawn
	RightSpawn dmqspecagent.UniversalSpawn

	LeftTTL            int32
	LeftMaxMessageSize uint64

	TopTTL            int32
	TopMaxMessageSize uint64

	RightTTL            int32
	RightMaxMessageSize uint64

	LogLeftToTopCommunication bool
	LogTopToLeftCommunication bool

	LogTopToRightCommunication bool
	LogRightToTopCommunication bool

	LogLeftLogs  bool
	LogTopLogs   bool
	LogRightLogs bool

	DisableAllLogs bool
}

type LoopTopoTestBenchTestFrameworkIntegration struct {
	LoopTopoTestBenchConfig

	AgentLogger          AgentLogger
	RecordingsComparator dmqspecagent.RecordingsComparator
	RecordingsLogger     RecordingsLogger
	BenchLogger          BenchLogger
}

type LoopTopoTestBench struct {
	config LoopTopoTestBenchTestFrameworkIntegration

	Left  dmqspecagent.Agent
	Top   dmqspecagent.Agent
	Right dmqspecagent.Agent

	leftTopSpy   *NodesCommunicationSpy
	topRightSpy  *NodesCommunicationSpy
	rightLeftSpy *NodesCommunicationSpy

	leftReady     chan struct{}
	leftConnected chan struct{}

	topReady     chan struct{}
	topConnected chan string

	rightReady     chan struct{}
	rightConnected chan struct{}

	alreadyStarted bool
	alreadyStopped bool
}

func NewLoopTopoTestBench(config LoopTopoTestBenchTestFrameworkIntegration) *LoopTopoTestBench {
	validateLoopTopoTestBenchConfig(config)

	return &LoopTopoTestBench{
		config: config,

		Left:  nil,
		Top:   nil,
		Right: nil,

		leftTopSpy:   nil,
		topRightSpy:  nil,
		rightLeftSpy: nil,

		leftReady:     make(chan struct{}),
		leftConnected: make(chan struct{}),

		topReady:     make(chan struct{}),
		topConnected: make(chan string),

		rightReady:     make(chan struct{}),
		rightConnected: make(chan struct{}),

		alreadyStarted: false,
		alreadyStopped: false,
	}
}

func NewGinkgoLoopTopoTestBench(config LoopTopoTestBenchConfig) *LoopTopoTestBench {
	return NewLoopTopoTestBench(LoopTopoTestBenchTestFrameworkIntegration{
		LoopTopoTestBenchConfig: config,

		BenchLogger:          GinkgoBenchLogger,
		AgentLogger:          GinkgoAgentLogger,
		RecordingsComparator: GomegaRecordingsComparator,
		RecordingsLogger:     GinkgoRecordingsLogger,
	})
}

func (bench *LoopTopoTestBench) Start() {
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

func (bench *LoopTopoTestBench) Stop(reason string) {
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

func (bench *LoopTopoTestBench) startSpies() {
	bench.benchLog("starting left-top spy")
	go bench.leftTopSpy.Start()

	bench.benchLog("starting top-right spy")
	go bench.topRightSpy.Start()

	bench.benchLog("starting right-left spy")
	go bench.rightLeftSpy.Start()
}

func (bench *LoopTopoTestBench) startAgents() {
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

	bench.benchLog("waiting for top agent to be ready")
	<-bench.topReady

	bench.benchLog("waiting for right agent to be ready")
	<-bench.rightReady
}

func (bench *LoopTopoTestBench) connectAgents() {
	bench.benchLog("listening on right agent: " + bench.topRightSpy.ToURL.String())
	bench.Right.Listen(dmqspecagent.ListenCommand{
		Address: bench.topRightSpy.ToURL.String(),
	})

	bench.benchLog("listening on top agent: " + bench.leftTopSpy.ToURL.String())
	bench.Top.Listen(dmqspecagent.ListenCommand{
		Address: bench.leftTopSpy.ToURL.String(),
	})

	bench.benchLog("listening on left agent: " + bench.rightLeftSpy.ToURL.String())
	bench.Left.Listen(dmqspecagent.ListenCommand{
		Address: bench.rightLeftSpy.ToURL.String(),
	})

	bench.benchLog("connecting from left agent to top agent")
	bench.Left.Connect(dmqspecagent.ConnectCommand{
		Address: bench.leftTopSpy.FromURL.String(),
	})

	bench.benchLog("waiting for left agent to connect to top agent")
	<-bench.leftConnected

	bench.benchLog("waiting for top agent to connect to left agent")
	<-bench.topConnected

	bench.benchLog("connecting from top agent to right agent")
	bench.Top.Connect(dmqspecagent.ConnectCommand{
		Address: bench.topRightSpy.FromURL.String(),
	})

	bench.benchLog("waiting for top agent to connect to right agent")
	<-bench.topConnected

	bench.benchLog("waiting for right agent to connect to top agent")
	<-bench.rightConnected

	bench.benchLog("connecting from right agent to left agent")
	bench.Right.Connect(dmqspecagent.ConnectCommand{
		Address: bench.rightLeftSpy.FromURL.String(),
	})

	bench.benchLog("waiting for right agent to connect to left agent")
	<-bench.rightConnected

	bench.benchLog("waiting for left agent to connect to right agent")
	<-bench.leftConnected
}

func (bench *LoopTopoTestBench) forceCleanUpTestBenchOnPanic() {
	if r := recover(); r != nil {
		bench.benchLog("force cleaning up test bench due to panic")

		bench.killAgents()
		bench.stopSpies()
	}
}

func (bench *LoopTopoTestBench) stopAgents(reason string) {
	bench.benchLog("stopping right agent")
	bench.Right.Stop(dmqspecagent.StopCommand{
		Reason: reason,
	})

	bench.benchLog("stopping top agent")
	bench.Top.Stop(dmqspecagent.StopCommand{
		Reason: reason,
	})

	bench.benchLog("stopping left agent")
	bench.Left.Stop(dmqspecagent.StopCommand{
		Reason: reason,
	})
}

func (bench *LoopTopoTestBench) stopSpies() {
	bench.benchLog("stopping left-top spy")
	bench.leftTopSpy.Stop()

	bench.benchLog("stopping top-right spy")
	bench.topRightSpy.Stop()

	bench.benchLog("stopping right-left spy")
	bench.rightLeftSpy.Stop()
}

func (bench *LoopTopoTestBench) killAgents() {
	bench.benchLog("force killing left agent")
	bench.Left.Kill()

	bench.benchLog("force killing top agent")
	bench.Top.Kill()

	bench.benchLog("force killing right agent")
	bench.Right.Kill()
}

func (bench *LoopTopoTestBench) configure() {
	bench.configureLeft()
	bench.configureTop()
	bench.configureRight()
	bench.configureLeftTopSpy()
	bench.configureTopRightSpy()
	bench.configureRightLeftSpy()
}

func (bench *LoopTopoTestBench) configureLeft() {
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

func (bench *LoopTopoTestBench) configureTop() {
	bench.benchLog("configuring top agent")

	bench.Top = dmqspecagent.NewUniversalAgent(bench.config.TopSpawn.NodeID)

	bench.Top.OnReady(func(_ dmqspecagent.ReadyNotification) {
		bench.topReady <- struct{}{}
	})

	bench.Top.OnConnectionEstablished(func(connected dmqspecagent.ConnectionEstablishedNotification) {
		bench.topConnected <- connected.BridgedNodeId
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

func (bench *LoopTopoTestBench) configureRight() {
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

func (bench *LoopTopoTestBench) configureLeftTopSpy() {
	bench.benchLog("configuring left-top spy")

	bench.leftTopSpy = NewNodesCommunicationSpy(bench.Left, bench.Top, NodesCommunicationSpyConfig{
		LogAToB: bench.config.LogLeftToTopCommunication && !bench.config.DisableAllLogs,
		LogBToA: bench.config.LogTopToLeftCommunication && !bench.config.DisableAllLogs,

		RecordingsComparator: bench.config.RecordingsComparator,
		RecordingsLogger:     bench.config.RecordingsLogger,
	})
}

func (bench *LoopTopoTestBench) configureTopRightSpy() {
	bench.benchLog("configuring top-right spy")

	bench.topRightSpy = NewNodesCommunicationSpy(bench.Top, bench.Right, NodesCommunicationSpyConfig{
		LogAToB: bench.config.LogTopToRightCommunication && !bench.config.DisableAllLogs,
		LogBToA: bench.config.LogRightToTopCommunication && !bench.config.DisableAllLogs,

		RecordingsComparator: bench.config.RecordingsComparator,
		RecordingsLogger:     bench.config.RecordingsLogger,
	})
}

func (bench *LoopTopoTestBench) configureRightLeftSpy() {
	bench.benchLog("configuring right-left spy")

	bench.rightLeftSpy = NewNodesCommunicationSpy(bench.Right, bench.Left, NodesCommunicationSpyConfig{
		LogAToB: bench.config.LogRightToTopCommunication && !bench.config.DisableAllLogs,
		LogBToA: bench.config.LogTopToRightCommunication && !bench.config.DisableAllLogs,

		RecordingsComparator: bench.config.RecordingsComparator,
		RecordingsLogger:     bench.config.RecordingsLogger,
	})
}

func (bench *LoopTopoTestBench) benchLog(message string) {
	if !bench.config.DisableAllLogs {
		bench.config.BenchLogger(message)
	}
}

func validateLoopTopoTestBenchConfig(config LoopTopoTestBenchTestFrameworkIntegration) {
	if config.LeftSpawn.NodeID == config.TopSpawn.NodeID {
		panic("Left and top node IDs must be different")
	}

	if config.TopSpawn.NodeID == config.RightSpawn.NodeID {
		panic("Top and right node IDs must be different")
	}

	if config.LeftSpawn.NodeID == config.RightSpawn.NodeID {
		panic("Left and right node IDs must be different")
	}

	if config.LeftTTL <= 0 {
		panic("Left TTL must be greater than 0")
	}

	if config.TopTTL <= 0 {
		panic("Top TTL must be greater than 0")
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

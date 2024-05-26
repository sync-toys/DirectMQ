package testbench

import (
	"time"

	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
)

type PairTopoTestBenchConfig struct {
	MasterSpawn dmqspecagent.UniversalSpawn
	SalveSpawn  dmqspecagent.UniversalSpawn

	MasterTTL            int32
	MasterMaxMessageSize uint64

	SalveTTL            int32
	SalveMaxMessageSize uint64

	LogMasterToSalveCommunication bool
	LogSalveToMasterCommunication bool

	LogMasterLogs bool
	LogSalveLogs  bool

	DisableAllLogs bool
}

type PairTopoTestBenchTestFrameworkIntegration struct {
	PairTopoTestBenchConfig

	AgentLogger          AgentLogger
	RecordingsComparator dmqspecagent.RecordingsComparator
	RecordingsLogger     RecordingsLogger
	BenchLogger          BenchLogger
}

type PairTopoTestBench struct {
	config PairTopoTestBenchTestFrameworkIntegration

	Master dmqspecagent.Agent
	Salve  dmqspecagent.Agent

	masterSalveSpy *NodesCommunicationSpy

	masterReady     chan struct{}
	masterConnected chan struct{}

	salveReady     chan struct{}
	salveConnected chan struct{}

	alreadyStarted bool
	alreadyStopped bool

	onBenchUp func()
}

func NewPairTopoTestBench(config PairTopoTestBenchTestFrameworkIntegration) *PairTopoTestBench {
	validatePairTopoTestBenchConfig(config)

	return &PairTopoTestBench{
		config: config,

		Master: nil,
		Salve:  nil,

		masterSalveSpy: nil,

		masterReady:     make(chan struct{}),
		masterConnected: make(chan struct{}),

		salveReady:     make(chan struct{}),
		salveConnected: make(chan struct{}),

		alreadyStarted: false,
		alreadyStopped: false,
	}
}

func NewGinkgoPairTopoTestBench(config PairTopoTestBenchConfig) *PairTopoTestBench {
	return NewPairTopoTestBench(PairTopoTestBenchTestFrameworkIntegration{
		PairTopoTestBenchConfig: config,

		BenchLogger:          GinkgoBenchLogger,
		AgentLogger:          GinkgoAgentLogger,
		RecordingsComparator: GomegaRecordingsComparator,
		RecordingsLogger:     GinkgoRecordingsLogger,
	})
}

func (t *PairTopoTestBench) Start() {
	defer t.forceCleanUpTestBenchOnPanic()

	if t.alreadyStarted {
		panic("test bench already started")
	}

	t.benchLog("starting test bench")
	t.alreadyStarted = true

	t.configure()

	t.benchLog("starting spy")
	go t.masterSalveSpy.Start()

	t.benchLog("spawning master")
	t.Master.Run(t.config.MasterSpawn, dmqspecagent.SetupCommand{
		TTL:            t.config.MasterTTL,
		MaxMessageSize: t.config.MasterMaxMessageSize,
		NodeID:         t.config.MasterSpawn.NodeID,
	})

	t.benchLog("spawning salve")
	t.Salve.Run(t.config.SalveSpawn, dmqspecagent.SetupCommand{
		TTL:            t.config.SalveTTL,
		MaxMessageSize: t.config.SalveMaxMessageSize,
		NodeID:         t.config.SalveSpawn.NodeID,
	})

	t.benchLog("waiting for master to be ready")
	<-t.masterReady

	t.benchLog("waiting for salve to be ready")
	<-t.salveReady

	if t.onBenchUp != nil {
		t.benchLog("calling onBenchUp callback")
		t.onBenchUp()
	}

	t.benchLog("listening on master")
	t.Master.Listen(
		dmqspecagent.ListenCommand{
			Address: t.masterSalveSpy.ToURL.String(),
		},
	)

	t.benchLog("connecting from salve")
	t.Salve.Connect(
		dmqspecagent.ConnectCommand{
			Address: t.masterSalveSpy.FromURL.String(),
		},
	)

	t.benchLog("waiting for salve to connect")
	<-t.salveConnected

	t.benchLog("waiting for master to connect")
	<-t.masterConnected

	t.benchLog("bench ready")
}

func (t *PairTopoTestBench) Stop(reason string) {
	defer t.forceCleanUpTestBenchOnPanic()

	if t.alreadyStopped {
		t.benchLog("test bench already stopped, ignoring stop request")
		return
	}

	t.benchLog("stopping test bench")
	t.alreadyStopped = true

	timeout := 10 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	done := make(chan struct{})
	go func() {
		t.stopAllServices(reason)
		close(done)
	}()

	select {
	case <-timer.C:
		panic("test bench stop timeout")
	case <-done:
		return
	}
}

func (t *PairTopoTestBench) Stopped() bool {
	return t.alreadyStopped
}

func (t *PairTopoTestBench) OnBenchUp(callback func()) {
	t.onBenchUp = callback
}

func (t *PairTopoTestBench) stopAllServices(reason string) {
	t.benchLog("gracefully stopping master agent")
	t.Master.Stop(dmqspecagent.StopCommand{Reason: reason})

	t.benchLog("gracefully stopping salve agent")
	t.Salve.Stop(dmqspecagent.StopCommand{Reason: reason})

	t.benchLog("stopping spy")
	t.masterSalveSpy.Stop()

	t.benchLog("test bench stopped")
}

func (t *PairTopoTestBench) SnapshotRecording(recordingID string) {
	t.benchLog("snapshotting recording as: " + recordingID)
	t.masterSalveSpy.SnapshotRecording(recordingID)
}

func (t *PairTopoTestBench) configure() {
	t.configureMaster()
	t.configureSalve()
	t.configureSpy()
}

func (t *PairTopoTestBench) configureMaster() {
	t.benchLog("configuring master")

	t.Master = dmqspecagent.NewUniversalAgent(t.config.MasterSpawn.NodeID)

	t.Master.OnReady(func(_ dmqspecagent.ReadyNotification) {
		t.masterReady <- struct{}{}
	})

	t.Master.OnConnectionEstablished(func(_ dmqspecagent.ConnectionEstablishedNotification) {
		t.masterConnected <- struct{}{}
	})

	t.Master.OnLogEntry(func(log string) {
		if !t.config.DisableAllLogs && t.config.LogMasterLogs {
			t.config.AgentLogger(log, t.Master.GetNodeID())
		}
	})

	t.Master.OnFatal(func(fatal dmqspecagent.FatalNotification) {
		t.config.AgentLogger("FATAL: "+fatal.Err, t.Master.GetNodeID())
	})
}

func (t *PairTopoTestBench) configureSalve() {
	t.benchLog("configuring salve")

	t.Salve = dmqspecagent.NewUniversalAgent(t.config.SalveSpawn.NodeID)

	t.Salve.OnReady(func(_ dmqspecagent.ReadyNotification) {
		t.salveReady <- struct{}{}
	})

	t.Salve.OnConnectionEstablished(func(_ dmqspecagent.ConnectionEstablishedNotification) {
		t.salveConnected <- struct{}{}
	})

	t.Salve.OnLogEntry(func(log string) {
		if !t.config.DisableAllLogs && t.config.LogSalveLogs {
			t.config.AgentLogger(log, t.Salve.GetNodeID())
		}
	})

	t.Salve.OnFatal(func(fatal dmqspecagent.FatalNotification) {
		t.config.AgentLogger("FATAL: "+fatal.Err, t.Master.GetNodeID())
	})
}

func (t *PairTopoTestBench) configureSpy() {
	t.benchLog("configuring spy")

	t.masterSalveSpy = NewNodesCommunicationSpy(t.Master, t.Salve, NodesCommunicationSpyConfig{
		LogAToB: !t.config.DisableAllLogs && t.config.LogMasterLogs,
		LogBToA: !t.config.DisableAllLogs && t.config.LogSalveLogs,

		RecordingsComparator: t.config.RecordingsComparator,
		RecordingsLogger:     t.config.RecordingsLogger,
	})
}

func (t *PairTopoTestBench) forceCleanUpTestBenchOnPanic() {
	r := recover()
	if r == nil {
		return
	}

	t.benchLog("test panicked, cleaning up test bench")

	if t.Master != nil {
		t.benchLog("force killing master")
		t.Master.Kill()
	}

	if t.Salve != nil {
		t.benchLog("force killing salve")
		t.Salve.Kill()
	}

	if t.masterSalveSpy != nil {
		t.benchLog("stopping spy")
		t.masterSalveSpy.Stop()
	}

	t.benchLog("test bench cleaned up, forwarding panic: " + r.(string))
	panic(r)
}

func (t *PairTopoTestBench) benchLog(log string) {
	if !t.config.DisableAllLogs {
		t.config.BenchLogger(log)
	}
}

func validatePairTopoTestBenchConfig(config PairTopoTestBenchTestFrameworkIntegration) {
	if config.MasterSpawn.NodeID == config.SalveSpawn.NodeID {
		panic("Master and salve node IDs must be different")
	}

	if config.MasterTTL <= 0 {
		panic("Master TTL must be greater than 0")
	}

	if config.SalveTTL <= 0 {
		panic("Salve TTL must be greater than 0")
	}

	if config.MasterMaxMessageSize < 0 {
		panic("Master max message size must be greater than or equal to 0")
	}

	if config.SalveMaxMessageSize < 0 {
		panic("Salve max message size must be greater than or equal to 0")
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

	if config.BenchLogger == nil {
		panic("Bench logger must be provided")
	}
}

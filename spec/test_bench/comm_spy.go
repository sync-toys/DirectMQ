package testbench

import (
	"net/url"
	"strconv"

	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
)

type RecordingsLogger func(message []byte, fromNode, toNode string)

type NodesCommunicationSpyConfig struct {
	LogAToB bool
	LogBToA bool

	RecordingsComparator dmqspecagent.RecordingsComparator
	RecordingsLogger     RecordingsLogger
}

type NodesCommunicationSpy struct {
	aAgent dmqspecagent.Agent
	bAgent dmqspecagent.Agent

	forwarder dmqspecagent.Forwarder
	recorder  *dmqspecagent.Recorder

	config NodesCommunicationSpyConfig

	FromURL *url.URL
	ToURL   *url.URL
}

func NewNodesCommunicationSpy(a dmqspecagent.Agent, b dmqspecagent.Agent, config NodesCommunicationSpyConfig) *NodesCommunicationSpy {
	fromURL := &url.URL{
		Scheme: "ws",
		Host:   "localhost:" + strconv.Itoa(getFreeTestingPort()),
		Path:   "/",
	}

	toURL := &url.URL{
		Scheme: "ws",
		Host:   "localhost:" + strconv.Itoa(getFreeTestingPort()),
		Path:   "/",
	}

	spy := &NodesCommunicationSpy{
		aAgent: a,
		bAgent: b,

		config: config,

		FromURL: fromURL,
		ToURL:   toURL,
	}

	spy.forwarder = dmqspecagent.NewWebsocketForwarder(fromURL, spy.bAgent.GetNodeID(), toURL, spy.aAgent.GetNodeID())
	spy.recorder = dmqspecagent.NewRecorder(spy.config.RecordingsComparator)

	spy.forwarder.OnMessage(spy.handleMessage)

	return spy
}

func (spy *NodesCommunicationSpy) handleMessage(message dmqspecagent.ForwardedMessage) {
	spy.recorder.Record(message.FromAlias, message.ToAlias, message.Message)

	if spy.config.LogAToB && message.FromAlias == spy.aAgent.GetNodeID() {
		spy.config.RecordingsLogger(message.Message, spy.aAgent.GetNodeID(), spy.bAgent.GetNodeID())
	}

	if spy.config.LogBToA && message.FromAlias == spy.bAgent.GetNodeID() {
		spy.config.RecordingsLogger(message.Message, spy.bAgent.GetNodeID(), spy.aAgent.GetNodeID())
	}
}

func (spy *NodesCommunicationSpy) Start() {
	spy.forwarder.StartWebsocketForwarder()
}

func (spy *NodesCommunicationSpy) Stop() {
	spy.forwarder.Abort()
}

func (spy *NodesCommunicationSpy) SnapshotRecording(recordingID string) {
	spy.recorder.SnapshotRecording(recordingID)
}

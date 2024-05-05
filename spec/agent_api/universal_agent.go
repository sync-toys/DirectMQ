package dmqspecagent

import (
	"bufio"
	"encoding/json"
	"io"
	"os/exec"
	"strings"

	directmq "github.com/sync-toys/DirectMQ/sdk/go"
)

type UniversalAgent struct {
	nodeID string

	cmd *exec.Cmd

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser

	// Agent API handlers
	onConnectionEstablished func(ConnectionEstablishedNotification)
	onConnectionLost        func(ConnectionLostNotification)
	onReady                 func(ReadyNotification)
	onFatal                 func(FatalNotification)
	onMessageReceived       func(MessageReceivedNotification)
	onLogEntry              func(string)

	// Diagnostics API handlers
	onPublication        func(OnPublicationNotification)
	onSubscription       func(OnSubscriptionNotification)
	onUnsubscribe        func(OnUnsubscribeNotification)
	onNetworkTermination func(OnNetworkTerminationNotification)

	subscribed chan SubscribedNotification
	stopped    chan StoppedNotification
}

var _ Agent = (*UniversalAgent)(nil)

func NewUniversalAgent() *UniversalAgent {
	return &UniversalAgent{
		subscribed: make(chan SubscribedNotification, 1),
		stopped:    make(chan StoppedNotification, 1),
	}
}

func (ua *UniversalAgent) spawn(executable UniversalSpawn) error {
	ua.cmd = exec.Command(executable.ExecutablePath, executable.Arguments...)

	stdin, err := ua.cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := ua.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := ua.cmd.StderrPipe()
	if err != nil {
		return err
	}

	ua.stdin = stdin
	ua.stdout = stdout
	ua.stderr = stderr

	if err := ua.cmd.Start(); err != nil {
		return err
	}

	go func() {
		ua.cmd.Wait()
		ua.cmd = nil
	}()

	return nil
}

func (ua *UniversalAgent) startNotificationLoop() {
	reader := bufio.NewReader(ua.stdout)

	for {
		message, err := reader.ReadString('\n')
		if err != nil && (err == io.EOF || strings.Contains(err.Error(), "file already closed")) {
			return
		}

		if err != nil {
			ua.onFatal(FatalNotification{Err: err.Error()})
			return
		}

		messageWithoutNewline := message[:len(message)-1]

		var notification UniversalNotification

		if err := json.Unmarshal([]byte(messageWithoutNewline), &notification); err != nil && ua.onLogEntry != nil {
			// Log messages printed to stdout are not valid JSON, they are just logs
			ua.onLogEntry(messageWithoutNewline)
			continue
		}

		ua.handleIncomingNotification(notification)
	}
}

func (ua *UniversalAgent) startErrorLoop() {
	reader := bufio.NewReader(ua.stderr)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		ua.onLogEntry(message[:len(message)-1])
	}
}

func (ua *UniversalAgent) handleIncomingNotification(n UniversalNotification) {
	if n.ConnectionEstablished != nil && ua.onConnectionEstablished != nil {
		ua.onConnectionEstablished(*n.ConnectionEstablished)
	}

	if n.ConnectionLost != nil && ua.onConnectionLost != nil {
		ua.onConnectionLost(*n.ConnectionLost)
	}

	if n.Ready != nil && ua.onReady != nil {
		ua.onReady(*n.Ready)
	}

	if n.Fatal != nil && ua.onFatal != nil {
		ua.onFatal(*n.Fatal)
	}

	if n.MessageReceived != nil && ua.onMessageReceived != nil {
		ua.onMessageReceived(*n.MessageReceived)
	}

	if n.OnPublication != nil && ua.onPublication != nil {
		ua.onPublication(*n.OnPublication)
	}

	if n.OnSubscription != nil && ua.onSubscription != nil {
		ua.onSubscription(*n.OnSubscription)
	}

	if n.OnUnsubscribe != nil && ua.onUnsubscribe != nil {
		ua.onUnsubscribe(*n.OnUnsubscribe)
	}

	if n.OnNetworkTermination != nil && ua.onNetworkTermination != nil {
		ua.onNetworkTermination(*n.OnNetworkTermination)
	}

	if n.Subscribed != nil {
		ua.subscribed <- *n.Subscribed
	}

	if n.Stopped != nil {
		ua.stopped <- *n.Stopped
	}
}

func (ua *UniversalAgent) write(cmd interface{}) {
	if ua.cmd == nil {
		panic("Agent process not available")
	}

	if ua.cmd.ProcessState != nil && ua.cmd.ProcessState.Exited() {
		panic("Agent process has exited")
	}

	cmdBytes, err := json.Marshal(cmd)
	if err != nil {
		panic("Failed to marshal command: " + err.Error())
	}

	cmdBytes = append(cmdBytes, '\n')

	_, err = ua.stdin.Write(cmdBytes)
	if err != nil {
		panic("Failed to write command: " + err.Error())
	}
}

////////////////////
// Connection API //
////////////////////

func (ua *UniversalAgent) Listen(cmd ListenCommand, exec UniversalSpawn) error {
	ua.nodeID = cmd.AsClientId

	if err := ua.spawn(exec); err != nil {
		return err
	}

	go ua.startErrorLoop()
	go ua.startNotificationLoop()

	ua.write(UniversalCommand{Listen: &cmd})

	return nil
}

func (ua *UniversalAgent) Connect(cmd ConnectCommand, exec UniversalSpawn) error {
	ua.nodeID = cmd.AsClientId

	if err := ua.spawn(exec); err != nil {
		return err
	}

	go ua.startErrorLoop()
	go ua.startNotificationLoop()

	ua.write(UniversalCommand{Connect: &cmd})

	return nil
}

func (ua *UniversalAgent) Stop(cmd StopCommand) {
	ua.write(UniversalCommand{Stop: &cmd})
	<-ua.stopped
	ua.Kill()
}

func (ua *UniversalAgent) Kill() {
	if ua.cmd == nil {
		return
	}

	ua.cmd.Process.Kill()
	ua.cmd = nil
}

///////////////
// Agent API //
///////////////

func (ua *UniversalAgent) OnReady(f func(ReadyNotification)) {
	ua.onReady = f
}

func (ua *UniversalAgent) OnFatal(f func(FatalNotification)) {
	ua.onFatal = f
}

func (ua *UniversalAgent) OnLogEntry(f func(string)) {
	ua.onLogEntry = f
}

func (ua *UniversalAgent) OnMessageReceived(f func(MessageReceivedNotification)) {
	ua.onMessageReceived = f
}

////////////////
// Native API //
////////////////

func (ua *UniversalAgent) Publish(cmd PublishCommand) {
	ua.write(UniversalCommand{Publish: &cmd})
}

func (ua *UniversalAgent) Subscribe(cmd SubscribeTopicCommand) directmq.SubscriptionID {
	ua.write(UniversalCommand{SubscribeTopic: &cmd})
	result := <-ua.subscribed
	return result.SubscriptionID
}

func (ua *UniversalAgent) Unsubscribe(cmd UnsubscribeTopicCommand) {
	ua.write(UniversalCommand{UnsubscribeTopic: &cmd})
}

/////////////////////
// Diagnostics API //
/////////////////////

func (ua *UniversalAgent) OnConnectionEstablished(f func(ConnectionEstablishedNotification)) {
	ua.onConnectionEstablished = f
}

func (ua *UniversalAgent) OnConnectionLost(f func(ConnectionLostNotification)) {
	ua.onConnectionLost = f
}

func (ua *UniversalAgent) OnPublication(f func(OnPublicationNotification)) {
	ua.onPublication = f
}

func (ua *UniversalAgent) OnSubscription(f func(OnSubscriptionNotification)) {
	ua.onSubscription = f
}

func (ua *UniversalAgent) OnUnsubscribe(f func(OnUnsubscribeNotification)) {
	ua.onUnsubscribe = f
}

func (ua *UniversalAgent) OnNetworkTermination(f func(OnNetworkTerminationNotification)) {
	ua.onNetworkTermination = f
}

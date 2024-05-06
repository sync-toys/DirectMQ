package testbench

import (
	"io"
	"reflect"

	"github.com/fatih/color"
	"github.com/onsi/gomega"
	dmqspec "github.com/sync-toys/DirectMQ/spec"
	dmqspecagent "github.com/sync-toys/DirectMQ/spec/agent_api"
)

func fontColor(colorCode color.Attribute) func(string, io.Writer) {
	return func(message string, writer io.Writer) {
		color.New(colorCode).Fprintln(writer, message)
	}
}

var GomegaRecordingsComparator dmqspecagent.RecordingsComparator = func(
	old []dmqspecagent.ForwardedMessageSnapRecord,
	new []dmqspecagent.ForwardedMessageSnapRecord,
	recordingID string,
	writeComparisonSnapshot func(),
) {
	message :=
		"Recording \"" + recordingID + "\" communication recording mismatch.\n\n" +
			"Nodes communication recordings are not equal.\n" +
			"Written snapshot for comparison next to the old snapshot file.\n" +
			"Please check the differences in the files and solve the issue.\n" +
			"If new recording is correct, please remove old snapshot and comparison file.\n" +
			"New recording will be saved as a current snapshot, you should commit it.\n"

	if !reflect.DeepEqual(old, new) {
		writeComparisonSnapshot()
	}

	gomega.Expect(old).To(gomega.Equal(new), message)
}

var GinkgoRecordingsLogger RecordingsLogger = func(message []byte, fromNode, toNode string) {
	formatter := fontColor(color.FgHiYellow)
	logComm := dmqspec.CreateLogger("comm", formatter)
	logComm("%s->%s: %s", fromNode, toNode, string(message))
}

var GinkgoAgentLogger AgentLogger = func(log string, nodeID string) {
	formatter := fontColor(color.FgHiBlue)
	logAgent := dmqspec.CreateLogger(nodeID, formatter)
	logAgent(log)
}

var GinkgoBenchLogger BenchLogger = func(log string) {
	formatter := fontColor(color.FgCyan)
	logBench := dmqspec.CreateLogger("bench", formatter)
	logBench(log)
}

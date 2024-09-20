package dmqspecagent

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sync-toys/DirectMQ/sdk/go/protocol"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func DecodeBinaryToJSON(payload []byte) interface{} {
	// decode from binary protobuf
	decodedBinary := new(protocol.DataFrame)
	if err := proto.Unmarshal(payload, decodedBinary); err != nil {
		panic("unable to decode message: " + err.Error())
	}

	// encode to json protobuf
	encodedJSON, err := protojson.Marshal(decodedBinary)
	if err != nil {
		panic("unable to encode message: " + err.Error())
	}

	// decode from json to interface
	var decodedJSON interface{}
	if err := json.Unmarshal(encodedJSON, &decodedJSON); err != nil {
		panic("unable to decode message: " + err.Error())
	}

	return decodedJSON
}

func DecodeBinaryToJSONString(payload []byte) string {
	decodedJSON := DecodeBinaryToJSON(payload)

	encodedJSON, err := json.Marshal(decodedJSON)
	if err != nil {
		panic("unable to encode message: " + err.Error())
	}

	return string(encodedJSON)
}

type ForwardedMessageSnapRecord struct {
	From    string
	To      string
	Message interface{}
}

type RecordingsComparator func(
	old []ForwardedMessageSnapRecord,
	new []ForwardedMessageSnapRecord,
	recordingID string,
	writeComparisonSnapshot func(),
)

type Recorder struct {
	recording []ForwardedMessageSnapRecord

	compareRecordings RecordingsComparator
}

func NewRecorder(
	compareRecordings RecordingsComparator,
) *Recorder {
	return &Recorder{
		recording:         []ForwardedMessageSnapRecord{},
		compareRecordings: compareRecordings,
	}
}

func (r *Recorder) Record(from string, to string, payload []byte) {
	r.recording = append(r.recording, ForwardedMessageSnapRecord{
		From:    from,
		To:      to,
		Message: DecodeBinaryToJSON(payload),
	})
}

func (r *Recorder) GetRecording() []ForwardedMessageSnapRecord {
	return r.recording
}

func (r *Recorder) SnapshotRecording(id string) {
	snapshotFile := "./snap/" + id + ".snap.json"

	if !r.snapshotExists(snapshotFile) {
		r.saveRecording(snapshotFile)
		return
	}

	oldRecording := r.loadRecording(snapshotFile)

	writeComparisonSnapshot := func() {
		comparisonFile := "./snap/" + id + ".CHANGED.snap.json"
		r.saveRecording(comparisonFile)
	}

	r.compareRecordings(oldRecording, r.recording, id, writeComparisonSnapshot)

	// we are not overriding the snapshot file if the recordings are the same or do not match
}

func (r *Recorder) snapshotExists(snapshotFile string) bool {
	_, err := os.Stat(snapshotFile)
	return !os.IsNotExist(err)
}

func (r *Recorder) saveRecording(snapshotFile string) {
	r.removeRecording(snapshotFile)

	r.mkdirs(snapshotFile)

	freshSnapshotFile, err := os.Create(snapshotFile)
	if err != nil {
		panic("unable to create snapshot file: " + err.Error())
	}

	defer freshSnapshotFile.Close()

	encoder := json.NewEncoder(freshSnapshotFile)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(r.recording)
	if err != nil {
		panic("unable to write snapshot file: " + err.Error())
	}
}

func (r *Recorder) loadRecording(snapshotFile string) []ForwardedMessageSnapRecord {
	oldSnapshotFile, err := os.Open(snapshotFile)
	if err != nil {
		panic("unable to open snapshot file: " + err.Error())
	}

	defer oldSnapshotFile.Close()

	var oldRecording []ForwardedMessageSnapRecord
	if err := json.NewDecoder(oldSnapshotFile).Decode(&oldRecording); err != nil {
		panic("unable to read snapshot file: " + err.Error())
	}

	return oldRecording
}

func (r *Recorder) mkdirs(snapshotFile string) {
	dir := filepath.Dir(snapshotFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic("unable to create snapshot directory: " + err.Error())
	}
}

func (r *Recorder) removeRecording(snapshotFile string) {
	if _, err := os.Stat(snapshotFile); os.IsNotExist(err) {
		return
	}

	if err := os.Remove(snapshotFile); err != nil {
		panic("unable to remove snapshot file: " + err.Error())
	}
}

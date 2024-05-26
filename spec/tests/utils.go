package dmqtests

import (
	"io"
	"time"

	"github.com/fatih/color"
	"github.com/onsi/ginkgo/v2"
	dmqspec "github.com/sync-toys/DirectMQ/spec"
)

var DISABLE_TEST_LOGGING bool = false

var _writeLog func(string, ...interface{}) = dmqspec.CreateLogger("test", func(message string, writer io.Writer) {
	color.New(color.FgGreen).Fprintln(writer, message)
})

func log(message string, args ...interface{}) {
	if !DISABLE_TEST_LOGGING {
		_writeLog(message, args...)
	}
}

func receiveWithTimeout[T interface{}](timeout int, ch chan T) T {
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		ginkgo.Fail("timeout")
		panic("timeout")
	case result := <-ch:
		return result
	}
}

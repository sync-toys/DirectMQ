package dmqspec

import (
	"fmt"
	"io"

	"github.com/onsi/ginkgo/v2"
)

func CreateLogger(prefix string, formatter func(message string, writer io.Writer)) func(string, ...interface{}) {
	return func(message string, args ...interface{}) {
		formattedPrefix := fmt.Sprintf("%-10s| ", prefix)
		formattedMessage := fmt.Sprintf(formattedPrefix+message, args...)
		formatter(
			formattedMessage,
			ginkgo.GinkgoWriter,
		)
	}
}

func NO_FORMATTER(message string, writer io.Writer) {
	fmt.Fprintln(writer, message)
}

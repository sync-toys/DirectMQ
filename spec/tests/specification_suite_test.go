package dmqtests

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSpecification(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Specification Suite")
}

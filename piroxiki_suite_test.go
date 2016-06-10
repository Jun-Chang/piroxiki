package piroxiki

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPiroxiki(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Piroxiki Suite")
}

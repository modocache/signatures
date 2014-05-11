package signatures_test

import (
	"github.com/modocache/signatures/signatures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/modocache/gory"
	"testing"
)

func TestSignatures(t *testing.T) {
	gory.Define("signature", signatures.Signature{}, func(factory gory.Factory) {
		factory["FirstName"] = "Jane"
		factory["LastName"] = "Doe"
		factory["Age"] = 27
		factory["Message"] = "I wholeheartedly support this petition!"
		factory["Email"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("jane-doe-%d@example.com", n)
		})
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Signatures Suite")
}

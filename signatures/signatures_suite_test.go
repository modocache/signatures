package signatures_test

import (
	. "github.com/modocache/signatures/signatures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/modocache/gory"
	"testing"
)

/*
Ginkgo generated this function
to kick off our unit tests. Hook into it
to define factories.
*/
func TestSignatures(t *testing.T) {
	defineFactories()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Signatures Suite")
}

/*
Define two factories: one for a valid signature,
and one for an invalid one (too young).
*/
func defineFactories() {
	gory.Define("signature", Signature{},
		func(factory gory.Factory) {
			factory["FirstName"] = "Jane"
			factory["LastName"] = "Doe"
			factory["Age"] = 27
			factory["Message"] = "I agree!"
			factory["Email"] = gory.Sequence(
				func(n int) interface{} {
					return fmt.Sprintf("jane-doe-%d@example.com", n)
				})
		})

	gory.Define("signatureTooYoung", Signature{},
		func(factory gory.Factory) {
			factory["FirstName"] = "Joey"
			factory["LastName"] = "Invalid"
			factory["Age"] = 10
			factory["Email"] = "joey-invalid@example.com"
		})
}

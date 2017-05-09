package concourse_test

import (
	"github.com/starkandwayne/credhub-resource/concourse"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {
	sillyBytes := []byte{0xFE, 0xED, 0xDE, 0xAD, 0xBE, 0xEF}
	sillyBytesSha1 := "33bf00cb7a45258748f833a47230124fcc8fa3a4"

	It("presents the SHA1 as a string", func() {
		Expect(concourse.NewVersion(sillyBytes, "foo.example.com")).To(Equal(concourse.Version{
			CredentialsSha1: sillyBytesSha1,
			Server:          "foo.example.com",
		}))
	})
})

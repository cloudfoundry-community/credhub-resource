package concourse_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/starkandwayne/credhub-resource/concourse"
)

var _ = Describe("NewSource", func() {
	It("converts the config into a Source", func() {
		config := []byte(`{
			"source": {
				"server": "foo.example.com",
				"username": "foo-user",
				"password": "foo-pass",
				"skip_tls_validation": false
			}
		}`)

		source, err := concourse.NewSource(config, "")
		Expect(err).NotTo(HaveOccurred())

		Expect(source).To(Equal(concourse.Source{
			Server:            "foo.example.com",
			Username:          "foo-user",
			Password:          "foo-pass",
			SkipTLSValidation: false,
		}))
	})

	Context("when decoding fails", func() {
		It("errors", func() {
			reader := []byte("not-json")

			_, err := concourse.NewSource(reader, "")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when a required parameter is missing", func() {
		It("returns an error with each missing parameter", func() {
			config := []byte("{}")

			_, err := concourse.NewSource(config, "")
			Expect(err).To(HaveOccurred())

			Expect(err.Error()).To(ContainSubstring("server"))
			Expect(err.Error()).To(ContainSubstring("username"))
			Expect(err.Error()).To(ContainSubstring("password"))
		})
	})
})

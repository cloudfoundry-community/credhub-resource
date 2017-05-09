package credhub_test

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-community/credhub-resource/clientfakes"
	"github.com/cloudfoundry-community/credhub-resource/concourse"
	. "github.com/cloudfoundry-community/credhub-resource/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		client     *CredhubClient
		httpClient *clientfakes.FakeHttpClient
	)
	BeforeEach(func() {
		httpClient = new(clientfakes.FakeHttpClient)
		client = NewCredhubClient(httpClient)
	})

	Describe("Login", func() {
		BeforeEach(func() {
			httpClient.DoReturnsOnCall(0, &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
            "app":{"version":"my-version","name":"CredHub"},
            "auth-server":{"url":"https://example.com"}
          }`)),
			}, nil)

			httpClient.DoReturnsOnCall(1, &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
            "access_token":"2YotnFZFEjr1zCsicMWpAA",
            "refresh_token":"5235FZFEjr1zCsicMWpAA",
            "token_type":"bearer",
            "expires_in":3600
          }`)),
			}, nil)
		})

		Context("With valid concourse.Source", func() {
			It("fetches a auth token", func() {
				err := client.Login(concourse.Source{
					Server:            "foo.example.com",
					Username:          "foo-user",
					Password:          "foo-pass",
					SkipTLSValidation: true,
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(client.Config).To(Equal(config.Config{
					ApiURL:             "foo.example.com",
					AuthURL:            "https://example.com",
					AccessToken:        "2YotnFZFEjr1zCsicMWpAA",
					RefreshToken:       "5235FZFEjr1zCsicMWpAA",
					InsecureSkipVerify: true,
				}))
			})
		})
	})
})

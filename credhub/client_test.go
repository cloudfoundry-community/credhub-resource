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

var _ = Describe("CredhubClient", func() {
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
            "auth-server":{"url":"https://uaa.example.com"}
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
					Server:            "https://credhub.example.com",
					Username:          "foo-user",
					Password:          "foo-pass",
					SkipTLSValidation: true,
				})

				Expect(err).ToNot(HaveOccurred())

				Expect(httpClient.DoArgsForCall(0).URL.Scheme).To(Equal("https"))
				Expect(httpClient.DoArgsForCall(0).URL.Host).To(Equal("credhub.example.com"))
				Expect(httpClient.DoArgsForCall(0).URL.Path).To(Equal("/info"))

				Expect(httpClient.DoArgsForCall(1).URL.Scheme).To(Equal("https"))
				Expect(httpClient.DoArgsForCall(1).URL.Host).To(Equal("uaa.example.com"))
				Expect(httpClient.DoArgsForCall(1).URL.Path).To(Equal("/oauth/token/"))

				Expect(client.Config).To(Equal(config.Config{
					ApiURL:             "https://credhub.example.com",
					AuthURL:            "https://uaa.example.com",
					AccessToken:        "2YotnFZFEjr1zCsicMWpAA",
					RefreshToken:       "5235FZFEjr1zCsicMWpAA",
					InsecureSkipVerify: true,
				}))
			})
		})
	})
})

package credhub_test

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-community/credhub-resource/clientfakes"
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
		client.Config = config.Config{
			ApiURL:      "https://credhub.example.com",
			AccessToken: "2YotnFZFEjr1zCsicMWpAA",
		}
	})

	Describe("FindAllCredentialPaths", func() {
		BeforeEach(func() {
			httpClient.DoReturnsOnCall(0, &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
            "credentials": [
              {
                "version_created_at": "2016-09-06T23:26:58Z",
                "name": "deploy1/dan/id.key"
              },
              {
                "version_created_at": "2016-09-06T23:26:58Z",
                "name": "dan.password"
              }
            ]
          }`)),
			}, nil)
		})
		Context("Given a path", func() {
			It("fetches a credentials list", func() {
				_, err := client.FindAllCredentialPaths("/foo")
				Expect(err).ToNot(HaveOccurred())

				Expect(httpClient.DoArgsForCall(0).URL.Scheme).To(Equal("https"))
				Expect(httpClient.DoArgsForCall(0).URL.Host).To(Equal("credhub.example.com"))
				Expect(httpClient.DoArgsForCall(0).URL.Path).To(Equal("/api/v1/data"))
				Expect(httpClient.DoArgsForCall(0).URL.RawQuery).To(Equal("path=%2Ffoo"))
			})
		})
	})

})

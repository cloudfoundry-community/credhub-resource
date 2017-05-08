package check_test

import (
	"bytes"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/credhub-cli/client/clientfakes"
	"github.com/starkandwayne/credhub-resource/check"
	"github.com/starkandwayne/credhub-resource/concourse"
	"github.com/starkandwayne/credhub-resource/credhub"
)

var _ = Describe("CheckCommand", func() {
	var (
		checkCommand check.CheckCommand
		client       credhub.CredhubClient
		httpClient   *clientfakes.FakeHttpClient
	)

	BeforeEach(func() {
		httpClient = new(clientfakes.FakeHttpClient)
		client = credhub.NewCredhubClient(httpClient)
		checkCommand = check.NewCheckCommand(&client)
	})

	Describe("Run", func() {
		var checkRequest concourse.CheckRequest

		Context("When the manifest SHA does not match with the version provided", func() {
			BeforeEach(func() {
				checkRequest = concourse.CheckRequest{
					Source: concourse.Source{
						Server:   "foo.example.com",
						Username: "foo-user",
						Password: "foo-pass",
					},
					Version: concourse.Version{},
				}

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

			It("different order should lead to same SHA1", func() {
				httpClient.DoReturnsOnCall(2, &http.Response{
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

				checkResponse, err := checkCommand.Run(checkRequest)
				Expect(httpClient.DoCallCount()).To(Equal(3))
				Expect(err).ToNot(HaveOccurred())
				Expect(checkResponse).To(Equal([]concourse.Version{
					{
						ManifestSha1: "182d2f06b2e6368bd2f22269362351616f5406ed",
						Target:       "foo.example.com",
					},
				}))
			})

			It("returns the SHA1 of the credhub keys", func() {
				httpClient.DoReturnsOnCall(2, &http.Response{
					StatusCode: 200,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
            "credentials": [
              {
                "name": "dan.password",
                "version_created_at": "2016-09-06T23:26:58Z"
              },
              {
                "name": "deploy1/dan/id.key",
                "version_created_at": "2016-09-06T23:26:58Z"
              }
            ]
          }`)),
				}, nil)

				checkResponse, err := checkCommand.Run(checkRequest)
				Expect(httpClient.DoCallCount()).To(Equal(3))
				Expect(err).ToNot(HaveOccurred())
				Expect(checkResponse).To(Equal([]concourse.Version{
					{
						ManifestSha1: "182d2f06b2e6368bd2f22269362351616f5406ed",
						Target:       "foo.example.com",
					},
				}))
			})

			It("returns different SHA1 when version_created_at changes", func() {
				httpClient.DoReturnsOnCall(2, &http.Response{
					StatusCode: 200,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
            "credentials": [
              {
                "name": "dan.password",
                "version_created_at": "1990-09-06T23:26:58Z"
              },
              {
                "name": "deploy1/dan/id.key",
                "version_created_at": "2016-09-06T23:26:58Z"
              }
            ]
          }`)),
				}, nil)

				checkResponse, err := checkCommand.Run(checkRequest)
				Expect(httpClient.DoCallCount()).To(Equal(3))
				Expect(err).ToNot(HaveOccurred())
				Expect(checkResponse).To(Equal([]concourse.Version{
					{
						ManifestSha1: "be1a1662163ce757245b0c8251acc26cabbb9583",
						Target:       "foo.example.com",
					},
				}))
			})
		})
	})
})

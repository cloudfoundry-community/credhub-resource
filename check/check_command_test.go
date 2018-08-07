package check_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"

	"github.com/cloudfoundry-community/credhub-resource/check"
	"github.com/cloudfoundry-community/credhub-resource/concourse"
	"github.com/cloudfoundry-community/credhub-resource/credhub/credhubfakes"
)

var _ = Describe("CheckCommand", func() {
	var (
		checkCommand check.CheckCommand
		client       *credhubfakes.FakeCredHub
	)

	BeforeEach(func() {
		client = new(credhubfakes.FakeCredHub)
		checkCommand = check.NewCheckCommand(client)
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
						Path:     "/foo-path",
					},
					Version: concourse.Version{},
				}

				client.FindByPathReturns(
					credentials.FindResults{
						Credentials: []credentials.Base{
							{
								Name:             "dan.password",
								VersionCreatedAt: "2016-09-06T23:26:58Z",
							},
							{
								Name:             "deploy1/dan/id.key",
								VersionCreatedAt: "2016-09-06T23:26:58Z",
							},
						},
					},
					nil,
				)
			})

			It("passes the path", func() {
				_, err := checkCommand.Run(checkRequest)
				Expect(client.FindByPathCallCount()).To(Equal(1))
				Expect(err).ToNot(HaveOccurred())
				Expect(client.FindByPathArgsForCall(0)).To(Equal("/foo-path"))
			})

			It("returns the SHA1 of the credhub keys", func() {
				checkResponse, err := checkCommand.Run(checkRequest)
				Expect(client.FindByPathCallCount()).To(Equal(1))
				Expect(err).ToNot(HaveOccurred())
				Expect(checkResponse).To(Equal([]concourse.Version{
					{
						CredentialsSha1: "c60309d18c0a68d2ac792e5a8717ad1f4d0bb25a",
						Server:          "foo.example.com",
					},
				}))
			})

			It("different order should lead to same SHA1", func() {
				client.FindByPathReturns(
					credentials.FindResults{
						Credentials: []credentials.Base{
							{
								Name:             "deploy1/dan/id.key",
								VersionCreatedAt: "2016-09-06T23:26:58Z",
							},
							{
								Name:             "dan.password",
								VersionCreatedAt: "2016-09-06T23:26:58Z",
							},
						},
					},
					nil,
				)

				checkResponse, err := checkCommand.Run(checkRequest)
				Expect(client.FindByPathCallCount()).To(Equal(1))
				Expect(err).ToNot(HaveOccurred())
				Expect(checkResponse).To(Equal([]concourse.Version{
					{
						CredentialsSha1: "c60309d18c0a68d2ac792e5a8717ad1f4d0bb25a",
						Server:          "foo.example.com",
					},
				}))
			})
		})
	})
})

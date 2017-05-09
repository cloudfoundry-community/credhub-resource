package credhub

import (
	"github.com/cloudfoundry-community/credhub-resource/concourse"
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

type CredhubClient struct {
	Config     config.Config
	HttpClient client.HttpClient
}

func BaseConfig(source concourse.Source) config.Config {
	return config.Config{ApiURL: source.Server, InsecureSkipVerify: source.SkipTLSValidation}
}

func NewCredhubClient(httpClient client.HttpClient) *CredhubClient {
	return &CredhubClient{HttpClient: httpClient}
}

func (c *CredhubClient) Login(source concourse.Source) error {
	c.Config = BaseConfig(source)
	serverInfo := actions.NewInfo(c.HttpClient, c.Config)
	info, err := serverInfo.GetServerInfo()
	if err != nil {
		return err
	}
	c.Config.AuthURL = info.AuthServer.Url

	serverInfo = actions.NewInfo(c.HttpClient, c.Config)
	token, err := serverInfo.GetAuthToken(source.Username, source.Password)
	if err != nil {
		return err
	}
	c.Config.AccessToken = token.AccessToken
	c.Config.RefreshToken = token.RefreshToken

	return nil
}

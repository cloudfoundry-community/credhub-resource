package in

import (
	"errors"

	"github.com/starkandwayne/credhub-resource/concourse"
	"github.com/starkandwayne/credhub-resource/credhub"
)

type InCommand struct {
	client credhub.CredhubClient
}

type InResponse struct {
	Version concourse.Version `json:"version"`
}

func NewInCommand(client credhub.CredhubClient) InCommand {
	return InCommand{client: client}
}

func (c InCommand) Run(inRequest concourse.InRequest, targetDir string) (InResponse, error) {
	err := c.client.Login(inRequest.Source)
	if err != nil {
		return InResponse{}, err
	}

	credentials, err := c.client.FindAllCredentialPaths("")
	if err != nil {
		return InResponse{}, err
	}

	actualVersion := concourse.NewVersion(credentials, inRequest.Source.Server)

	if actualVersion.Server != inRequest.Version.Server {
		return InResponse{}, errors.New("Credhub server is different than configured source")
	}

	if actualVersion.CredentialsSha1 != inRequest.Version.CredentialsSha1 {
		return InResponse{}, errors.New("Credhub credentials fingerprint can used as input")
	}

	return InResponse{Version: actualVersion}, nil
}

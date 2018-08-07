package in

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/cloudfoundry-community/credhub-resource/concourse"
	"github.com/cloudfoundry-community/credhub-resource/credhub"
)

type InCommand struct {
	client credhub.CredHub
}

type InResponse struct {
	Version concourse.Version `json:"version"`
}

func NewInCommand(client credhub.CredHub) InCommand {
	return InCommand{client: client}
}

func (c InCommand) Run(inRequest concourse.InRequest, targetDir string) (InResponse, error) {
	credentials, err := c.client.FindByPath(inRequest.Source.Path)
	if err != nil {
		return InResponse{}, err
	}
	sort.Slice(credentials.Credentials, func(i, j int) bool {
		return credentials.Credentials[i].Name < credentials.Credentials[j].Name
	})
	raw, err := json.Marshal(credentials)
	if err != nil {
		return InResponse{}, err
	}

	actualVersion := concourse.NewVersion(raw, inRequest.Source.Server)

	if actualVersion.Server != inRequest.Version.Server {
		return InResponse{}, errors.New("Credhub server is different than configured source")
	}

	if actualVersion.CredentialsSha1 != inRequest.Version.CredentialsSha1 {
		return InResponse{}, errors.New("Credhub credentials fingerprint can used as input")
	}

	return InResponse{Version: actualVersion}, nil
}

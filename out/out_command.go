package out

import (
	"encoding/json"
	"sort"

	"github.com/cloudfoundry-community/credhub-resource/concourse"
	"github.com/cloudfoundry-community/credhub-resource/credhub"
)

type OutResponse struct {
	Version  concourse.Version    `json:"version"`
	Metadata []concourse.Metadata `json:"metadata"`
}

type OutCommand struct {
	client             credhub.CredHub
	resourcesDirectory string
}

func NewOutCommand(client credhub.CredHub, resourcesDirectory string) OutCommand {
	return OutCommand{
		client:             client,
		resourcesDirectory: resourcesDirectory,
	}
}

func (c OutCommand) Run(outRequest concourse.OutRequest) (OutResponse, error) {
	credentials, err := c.client.FindByPath(outRequest.Source.Path)
	if err != nil {
		return OutResponse{}, err
	}
	sort.Slice(credentials.Credentials, func(i, j int) bool {
		return credentials.Credentials[i].Name < credentials.Credentials[j].Name
	})
	raw, err := json.Marshal(credentials)
	if err != nil {
		return OutResponse{}, err
	}

	version := concourse.NewVersion(raw, outRequest.Source.Server)

	concourseOutput := OutResponse{
		Version:  version,
		Metadata: append([]concourse.Metadata{}),
	}

	return concourseOutput, nil
}

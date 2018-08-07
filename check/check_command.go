package check

import (
	"encoding/json"
	"sort"

	"github.com/cloudfoundry-community/credhub-resource/concourse"
	"github.com/cloudfoundry-community/credhub-resource/credhub"
)

type CheckCommand struct {
	client credhub.CredHub
}

func NewCheckCommand(client credhub.CredHub) CheckCommand {
	return CheckCommand{client: client}
}

func (c CheckCommand) Run(checkRequest concourse.CheckRequest) ([]concourse.Version, error) {
	credentials, err := c.client.FindByPath(checkRequest.Source.Path)
	if err != nil {
		return []concourse.Version{}, err
	}
	sort.Slice(credentials.Credentials, func(i, j int) bool {
		return credentials.Credentials[i].Name < credentials.Credentials[j].Name
	})
	raw, err := json.Marshal(credentials)
	if err != nil {
		return []concourse.Version{}, err
	}

	version := concourse.NewVersion(raw, checkRequest.Source.Server)

	var concourseOutput = []concourse.Version{}
	if version != checkRequest.Version {
		concourseOutput = append(concourseOutput, version)
	}

	return concourseOutput, nil
}

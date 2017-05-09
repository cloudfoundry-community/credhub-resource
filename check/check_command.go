package check

import (
	"github.com/starkandwayne/credhub-resource/concourse"
	"github.com/starkandwayne/credhub-resource/credhub"
)

type CheckCommand struct {
	client *credhub.CredhubClient
}

func NewCheckCommand(client *credhub.CredhubClient) CheckCommand {
	return CheckCommand{client: client}
}

func (c CheckCommand) Run(checkRequest concourse.CheckRequest) ([]concourse.Version, error) {
	err := c.client.Login(checkRequest.Source)
	if err != nil {
		return []concourse.Version{}, err
	}

	credentials, err := c.client.FindAllCredentialPaths("")
	if err != nil {
		return []concourse.Version{}, err
	}

	version := concourse.NewVersion([]byte(credentials), checkRequest.Source.Server)

	var concourseOutput = []concourse.Version{}
	if version != checkRequest.Version {
		concourseOutput = append(concourseOutput, version)
	}

	return concourseOutput, nil
}

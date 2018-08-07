package credhub

import (
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
)

type CredHub interface {
	FindByPath(string) (credentials.FindResults, error)
}

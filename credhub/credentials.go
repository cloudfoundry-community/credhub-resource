package credhub

import (
	"bytes"
	"sort"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
	"github.com/ugorji/go/codec"
)

func (c CredhubClient) FindAllCredentialPaths(path string) ([]byte, error) {
	action := actions.NewAction(
		repositories.NewCredentialQueryRepository(c.HttpClient),
		c.Config,
	)

	credentials, err := action.DoAction(
		client.NewFindCredentialsByPathRequest(c.Config, path), path)
	if err != nil {
		return []byte{}, err
	}

	canonicallizedJSON, err := canonicallizeCredentialsJSON(credentials.ToJson())
	if err != nil {
		return []byte{}, err
	}

	return []byte(canonicallizedJSON), nil
}

func canonicallizeCredentialsJSON(json string) (string, error) {
	var creds struct {
		Credentials []struct {
			Name             string    `json:"name"`
			VersionCreatedAt time.Time `json:"version_created_at"`
		} `json:"credentials"`
	}
	err := codec.NewDecoder(strings.NewReader(json), new(codec.JsonHandle)).Decode(&creds)
	if err != nil {
		return "", err
	}

	sort.Slice(creds.Credentials, func(i, j int) bool { return creds.Credentials[i].Name < creds.Credentials[j].Name })

	w := new(bytes.Buffer)
	h := new(codec.JsonHandle)
	h.Canonical = true
	enc := codec.NewEncoder(w, h)
	err = enc.Encode(creds)

	return string(w.Bytes()), nil
}

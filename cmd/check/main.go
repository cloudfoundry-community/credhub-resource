package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"

	"github.com/cloudfoundry-community/credhub-resource/check"
	"github.com/cloudfoundry-community/credhub-resource/concourse"
)

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read configuration: %s\n", err)
		os.Exit(1)
	}

	checkRequest, err := concourse.NewCheckRequest(stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid parameters: %s\n", err)
		os.Exit(1)
	}

	client, err := credhub.New(
		checkRequest.Source.Server,
		credhub.SkipTLSValidation(checkRequest.Source.SkipTLSValidation),
		credhub.Auth(auth.UaaClientCredentials(
			checkRequest.Source.Username,
			checkRequest.Source.Password,
		)),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating CredHub client: %s\n", err)
		os.Exit(1)
	}

	checkCommand := check.NewCheckCommand(client)
	checkResponse, err := checkCommand.Run(checkRequest)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	concourseOutputFormatted, err := json.MarshalIndent(checkResponse, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate version: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s", concourseOutputFormatted)
}

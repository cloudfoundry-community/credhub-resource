package main

import (
	"encoding/json"
	"fmt"
	"os"

	"io/ioutil"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/starkandwayne/credhub-resource/check"
	"github.com/starkandwayne/credhub-resource/concourse"
	"github.com/starkandwayne/credhub-resource/credhub"
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

	checkCommand := check.NewCheckCommand(credhub.NewCredhubClient(
		client.NewHttpClient(
			credhub.BaseConfig(checkRequest.Source))))
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

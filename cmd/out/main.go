package main

import (
	"encoding/json"
	"fmt"
	"os"

	"io/ioutil"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/starkandwayne/credhub-resource/concourse"
	"github.com/starkandwayne/credhub-resource/credhub"
	"github.com/starkandwayne/credhub-resource/out"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr,
			"not enough args - usage: %s <sources directory>\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	sourcesDir := os.Args[1]

	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read configuration: %s\n", err)
		os.Exit(1)
	}

	outRequest, err := concourse.NewOutRequest(stdin, sourcesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid parameters: %s\n", err)
		os.Exit(1)
	}

	outCommand := out.NewOutCommand(credhub.NewCredhubClient(
		client.NewHttpClient(credhub.BaseConfig(outRequest.Source))), sourcesDir)
	outResponse, err := outCommand.Run(outRequest)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	concourseOutputFormatted, err := json.MarshalIndent(outResponse, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate version: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s", concourseOutputFormatted)
}

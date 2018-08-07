package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"

	"github.com/cloudfoundry-community/credhub-resource/concourse"
	"github.com/cloudfoundry-community/credhub-resource/in"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr,
			"not enough args - usage: %s <target directory>\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read configuration: %s\n", err)
		os.Exit(1)
	}

	inRequest, err := concourse.NewInRequest(stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid parameters: %s\n", err)
		os.Exit(1)
	}

	targetDir := os.Args[1]

	if inRequest.Source.Server == concourse.MissingServer {
		inResponse := in.InResponse{Version: inRequest.Version}
		printResponse(inResponse)
		os.Exit(0)
	}

	client, err := credhub.New(
		inRequest.Source.Server,
		credhub.SkipTLSValidation(inRequest.Source.SkipTLSValidation),
		credhub.Auth(auth.UaaClientCredentials(
			inRequest.Source.Username,
			inRequest.Source.Password,
		)),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating CredHub client: %s\n", err)
		os.Exit(1)
	}

	inCommand := in.NewInCommand(client)
	inResponse, err := inCommand.Run(inRequest, targetDir)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	printResponse(inResponse)
}

func printResponse(inResponse in.InResponse) {
	concourseInputFormatted, err := json.MarshalIndent(inResponse, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not generate version: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s", concourseInputFormatted)
}

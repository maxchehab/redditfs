package main

import (
	"fmt"

	survey "gopkg.in/AlecAivazis/survey.v1"
)

// List displays all files and directories that are available
func List(args []string, _ []Command) error {
	manifest, err := RetrieveManifestFromReddit(testSubreddit)
	if err != nil {
		return err
	}

	var repositories []string
	for _, repo := range manifest.Repositories {
		repositories = append(repositories, repo.Name)
	}

	prompt := &survey.MultiSelect{
		Message: "Which repositories would you like to download? (You may select more than one)",
		Options: repositories,
	}
	response := []string{}
	survey.AskOne(prompt, &response, nil)

	if err != nil {
		return err
	}

	fmt.Printf("You choose %v\n", response)
	return err
}

package main

import (
	"os"
	"sync"

	survey "gopkg.in/AlecAivazis/survey.v1"
)

// List displays all files and directories that are available
func List(args []string, _ []Command) error {
	manifest, err := RetrieveManifestFromReddit(testSubreddit)

	if err != nil {
		return err
	}

	var options []string
	for _, repo := range manifest.Repositories {
		options = append(options, repo.Name)
	}

	prompt := &survey.MultiSelect{
		Message: "Select repositories to download",
		Options: options,
	}
	response := []string{}
	survey.AskOne(prompt, &response, nil)

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	repositories := manifest.FilterReposByNames(response)

	var wg sync.WaitGroup

	downloadResponses := make(chan error, len(repositories))

	for _, repository := range repositories {
		wg.Add(1)
		go func(repository Repository) {
			defer wg.Done()
			downloadResponses <- repository.Download(path)
		}(repository)
	}
	wg.Wait()
	close(downloadResponses)

	for downloadResponse := range downloadResponses {
		if downloadResponse != nil {
			return downloadResponse
		}
	}
	return err
}

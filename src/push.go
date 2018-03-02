package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/ttacon/chalk"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Push displays all files and directories that are available
func Push(args []string, _ []Command) error {
	session := GetSession()
	manifest, err := RetrieveManifestFromReddit(testSubreddit, session)

	if err != nil {
		return err
	}

	var options []string
	for _, repo := range manifest.Repositories {
		options = append(options, repo.Name)
	}

	selectedPath, err := os.Getwd()
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Confirm push of repository [%v%v%v]", chalk.Cyan, selectedPath, chalk.Reset)

	confirm := false
	prompt := &survey.Confirm{
		Message: message,
	}
	survey.AskOne(prompt, &confirm, nil)

	if !confirm {
		return nil
	}

	pathsToUpload, err := GetAllFilesInDirectory(selectedPath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	files := make(chan File, len(pathsToUpload))
	fileErrors := make(chan error, len(pathsToUpload))
	for _, path := range pathsToUpload {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			file, err := UploadFileByPath(path, selectedPath, session)
			files <- file
			fileErrors <- err
		}(path)
	}
	wg.Wait()
	close(files)
	close(fileErrors)

	for err := range fileErrors {
		if err != nil {
			return err
		}
	}
	return err
}

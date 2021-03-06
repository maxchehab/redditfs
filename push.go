package main

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/maxchehab/geddit"
	"github.com/ttacon/chalk"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Push displays all files and directories that are available
func Push(args []string, _ []Command) (err error) {
	username, password, subreddit := GetCredentials()
	session, err := GetSession(username, password)
	if err != nil {
		return err
	}

	manifest, err := RetrieveManifestFromReddit(subreddit, session)

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

	message := fmt.Sprintf("Confirm push of repository [%v%v%v]", chalk.Cyan, path.Base(selectedPath), chalk.Reset)

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
		func(path string) {
			defer wg.Done()
			file, err := UploadFileByPath(path, selectedPath, session, subreddit)
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

	var respository Repository
	respository.Name = path.Base(selectedPath)
	for file := range files {
		respository.Files = append(respository.Files, file)
	}

	if manifest.Contains(respository) {
		message := fmt.Sprintf("The repository [%v%v%v] has already been created. Would you like to update or rename the repository.", chalk.Cyan, respository.Name, chalk.Reset)
		action := ""
		prompt := &survey.Select{
			Message: message,
			Options: []string{"update", "rename"},
		}
		survey.AskOne(prompt, &action, nil)

		if action == "rename" {
			message = fmt.Sprintf("Rename [%v%v%v] to:", chalk.Cyan, respository.Name, chalk.Reset)
			rename := ""
			prompt := &survey.Input{
				Message: message,
			}
			survey.AskOne(prompt, &rename, nil)
			Prompt(fmt.Sprintf("Renamed [%v%v%v] to [%v%v%v]", chalk.Cyan, respository.Name, chalk.ResetColor, chalk.Cyan, rename, chalk.ResetColor))
			respository.Name = rename
			manifest.Repositories = append(manifest.Repositories, respository)
		} else {
			manifest.Update(respository)
		}
	} else {
		manifest.Repositories = append(manifest.Repositories, respository)
	}

	Prompt(fmt.Sprintf("Updating manifest for [%v%v%v]", chalk.Cyan, respository.Name, chalk.ResetColor))
	if len(manifest.Location) > 0 {
		session.EditUserText(geddit.NewEdit(manifest.ToString(), "t3_"+manifest.Location))
	} else {
		session.Submit(geddit.NewTextSubmission(subreddit, "manifest.json", manifest.ToString(), false, nil))
	}

	return
}

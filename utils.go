package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/maxchehab/geddit"
	gitignore "github.com/monochromegane/go-gitignore"
	"github.com/ttacon/chalk"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// FileExists checks fi a path to a file exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// GetAllFilesInDirectory searches a provided directory and returns a slice
// representing all files within selected directory.
func GetAllFilesInDirectory(path string) (paths []string, err error) {
	var ignoreMatcher gitignore.IgnoreMatcher
	gitIgnoreExists := FileExists(path + "/.redditfsignore")
	if gitIgnoreExists {
		ignoreMatcher, err = gitignore.NewGitIgnore(path + "/.redditfsignore")
		if err != nil {
			return paths, err
		}
	}

	err = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if gitIgnoreExists {
				if !ignoreMatcher.Match(path, false) {
					paths = append(paths, path)
				}
			} else {
				paths = append(paths, path)
			}
		} else if gitIgnoreExists && ignoreMatcher.Match(path, true) {
			return filepath.SkipDir
		}
		return nil
	})
	return paths, err
}

// Subfolders returns a slice of subfolders (recursive), including the folder provided.
func Subfolders(path string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			paths = append(paths, newPath)
		}
		return nil
	})
	return paths
}

// CreateNewOauthCredentials will prompt and write OAUTH credentials form the user
func CreateNewOauthCredentials() (name string, id string, secret string, read bool, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	credentialPath := usr.HomeDir + "/.redditfscredentials"

	Prompt("Please create your authorized reddit application.")
	Prompt("This information will be stored at " + credentialPath)
	Prompt("To learn more, go to https://github.com/maxchehab/redditfs")

	namePrompt := &survey.Input{
		Message: "Name:",
	}
	idPrompt := &survey.Input{
		Message: "ID:",
	}
	secretPrompt := &survey.Input{
		Message: "Secret:",
	}

	survey.AskOne(namePrompt, &name, nil)
	survey.AskOne(idPrompt, &id, nil)
	survey.AskOne(secretPrompt, &secret, nil)

	output := name + "\n" + id + "\n" + secret
	err = ioutil.WriteFile(credentialPath, []byte(output), 0600)
	read = true
	return
}

// GetOauthCredentials returns OAUTH credentials
func GetOauthCredentials() (name string, id string, secret string, read bool, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}

	credentialPath := usr.HomeDir + "/.redditfscredentials"

	if _, err := os.Stat(credentialPath); os.IsNotExist(err) {
		return CreateNewOauthCredentials()
	}

	if _, err := os.Stat(credentialPath); err == nil {
		file, err := os.Open(credentialPath)
		if err != nil {
			return name, id, secret, read, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		name = scanner.Text()
		scanner.Scan()
		id = scanner.Text()
		scanner.Scan()
		secret = scanner.Text()
		read = true
	}

	return
}

// GetCredentials reads username and password from stdinput
func GetCredentials() (username string, password string, subreddit string) {
	Prompt("Please enter your reddit username, password and subreddit.")
	usernamePrompt := &survey.Input{
		Message: "Username:",
	}
	passwordPrompt := &survey.Password{
		Message: "Password:",
	}
	subredditPrompt := &survey.Input{
		Message: "Subreddit:",
	}

	survey.AskOne(usernamePrompt, &username, nil)
	survey.AskOne(passwordPrompt, &password, nil)
	survey.AskOne(subredditPrompt, &subreddit, nil)

	return
}

//Prompt will display an interface on the screen
func Prompt(input interface{}) {
	fmt.Println(chalk.Bold.TextStyle(fmt.Sprintf("%v>%v %v", chalk.Green, chalk.ResetColor, input)))
}

// WriteByteStringToFile takes a string of bytes seperated by a space
// and writes the binary data to a file
func WriteByteStringToFile(input string, file string) (err error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		os.MkdirAll(path.Dir(file), os.ModePerm)
		os.Create(file)
	}
	output, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer output.Close()
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		s, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}
		b := byte(s)
		output.Write([]byte{b})
	}
	return
}

// UploadFileByPath uploads a file and returns a file object
func UploadFileByPath(absolutePath string, selectedPath string, session *geddit.OAuthSession, subreddit string) (file File, err error) {
	relativePath := absolutePath[len(selectedPath):]
	Prompt(fmt.Sprintf("Uploading %v%v%v", chalk.Cyan, relativePath, chalk.Cyan))
	// file.Path = path.Dir(relativePath) + "/"
	relativePathDir := path.Dir(relativePath)
	if relativePathDir[len(relativePathDir)-1:] == "/" {
		file.Path = relativePathDir
	} else {
		file.Path = relativePathDir + "/"
	}
	file.Name = path.Base(relativePath)
	buffer := make([]byte, 8192)

	input, err := os.Open(absolutePath)
	if err != nil {
		return file, err
	}

	for {
		_, err := io.ReadFull(input, buffer)
		if err == io.EOF {
			break
		} else {
			buffer = bytes.Trim(buffer, string(byte(0)))
			err := file.UploadBuffer(buffer, session, subreddit)
			if err != nil {
				return file, err
			}
		}
	}
	return
}

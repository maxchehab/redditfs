package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
		fmt.Println(gitIgnoreExists)
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

// Credentials reads username and password from stdinput
func Credentials() (username string, password string) {

	usernamePrompt := &survey.Input{
		Message: "Username:",
	}
	passwordPrompt := &survey.Password{
		Message: "Password:",
	}

	survey.AskOne(usernamePrompt, &username, nil)
	survey.AskOne(passwordPrompt, &password, nil)
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
func UploadFileByPath(absolutePath string, selectedPath string, session *geddit.OAuthSession) (file File, err error) {
	relativePath := absolutePath[len(selectedPath):]
	Prompt(fmt.Sprintf("Uploading %v%v%v", chalk.Cyan, relativePath, chalk.Cyan))
	file.Path = path.Dir(relativePath) + "/"
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
			err := file.UploadBuffer(buffer, session)
			if err != nil {
				return file, err
			}
		}
	}
	return
}

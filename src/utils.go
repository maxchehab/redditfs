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
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// GetAllFilesInDirectory searches a provided directory and returns a slice
// representing all files within selected directory.
func GetAllFilesInDirectory(path string) (paths []string, err error) {
	err = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			paths = append(paths, path)
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
		Message: "Pasword:",
	}

	survey.AskOne(usernamePrompt, &username, nil)
	survey.AskOne(passwordPrompt, &password, nil)
	return
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
func UploadFileByPath(path string, selectedPath string, session *geddit.OAuthSession) (file File, err error) {
	fmt.Printf("Uploading %v", path)
	buffer := make([]byte, 8192)
	input, err := os.Open(path)
	if err != nil {
		return file, err
	}

	for {
		_, err := io.ReadFull(input, buffer)
		if err == io.EOF {
			break
		} else {
			_, err := file.UploadBuffer(buffer, session)
			if err != nil {
				return file, err
			}
		}
	}
	return
}

// input, _ := os.Open("./test/input.jpg")

// output, _ := os.OpenFile("./test/middle.csv", os.O_APPEND|os.O_WRONLY, 0600)

// buf := make([]byte, 2048)
// for {
// 	_, err := io.ReadFull(input, buf)
// 	if err == io.EOF {
// 		break
// 	} else {
// 		for _, b := range buf {
// 			o := strconv.Itoa(int(b))
// 			io.WriteString(output, o+" ")
// 		}
// 	}
// }
// input.Close()
// output.Close()

// input, _ = os.Open("./test/middle.csv")
// output, _ = os.OpenFile("./test/output.jpg", os.O_APPEND|os.O_WRONLY, 0600)
// scanner := bufio.NewScanner(input)
// scanner.Split(bufio.ScanWords)
// for scanner.Scan() {
// 	s, err := strconv.Atoi(scanner.Text())
// 	if err != nil {
// 		panic(err)
// 	}
// 	b := byte(s)
// 	output.Write([]byte{b})
// }

// input.Close()
// output.Close()

package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

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
func Credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))

	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}

// WriteByteStringToFile takes a string of bytes seperated by a space
// and writes the binary data to a file
func WriteByteStringToFile(input string, file string) (err error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		os.MkdirAll(path.Dir(file), os.ModePerm)
		os.Create(file)
	}
	output, err := os.OpenFile(file, os.O_TRUNC|os.O_WRONLY, 0600)
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

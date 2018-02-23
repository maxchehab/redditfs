package main

import (
	"fmt"

	"github.com/ryanuber/columnize"
)

// Command is a stucture representing each command
type Command struct {
	Action func(args []string, commands []Command) error
	Name   string
	Tokens []string
	Help   string
}

// Execute will print the command's Help or execute the commands Action
func (command Command) Execute(args []string, commands []Command) error {
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Println(command.Help)
	} else {
		return command.Action(args, commands)
	}
	return nil
}

// Contains reports if an argument can be resolved as a token
func (command Command) Contains(arg string) bool {
	for _, token := range command.Tokens {
		if token == arg {
			return true
		}
	}
	return false
}

// List displays all files and directories that are available
func List(args []string, _ []Command) error {
	manifest, err := RetrieveManifestFromReddit(testSubreddit)

	for _, file := range manifest.Files {
		fmt.Println(file.Name)
	}

	return err
}

// Help displays all commands that are available
func Help(args []string, commands []Command) error {
	lines := []string{
		"Name \t Command \t Information",
	}

	for _, command := range commands {
		lines = append(lines, command.Name+" \t "+fmt.Sprint(command.Tokens)+" \t "+command.Help)
	}

	config := columnize.DefaultConfig()
	config.Delim = "\t"

	fmt.Println("Usage: redditfs [command] [args]")
	fmt.Println(columnize.Format(lines, config))

	return nil
}

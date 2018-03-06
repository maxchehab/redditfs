package main

import (
	"fmt"
	"os"

	"github.com/ttacon/chalk"
)

func main() {
	commands := []Command{
		Command{Help, "Help", []string{"--help", "-h", "help"}, "Displays help."},
		Command{Pull, "Pull", []string{"pull", "ls", "list"}, "Display, select, and download any repositories that are available."},
		Command{Push, "Push", []string{"push", "add"}, "Upload current directory."},
	}

	args := os.Args[1:]
	if len(args) > 0 {
		for _, command := range commands {
			if command.Contains(args[0]) {
				err := command.Execute(args[1:], commands)
				if err != nil {
					fmt.Println(chalk.Red, "I'm sorry, but this command didn't work as expected.", chalk.Reset)
					panic(err)
				}
				return
			}
		}
	}

	commands[0].Action([]string{}, commands)
}

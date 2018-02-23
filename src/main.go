package main

import "os"

func main() {
	commands := []Command{
		Command{Help, "Help", []string{"--help", "-h", "help"}, "Displays help."},
		Command{List, "List", []string{"ls", "list"}, "Display all files and directories that are available"},
	}

	args := os.Args[1:]
	if len(args) > 0 {
		for _, command := range commands {
			if command.Contains(args[0]) {
				command.Execute(args[1:], commands)
				return
			}
		}
	}

	commands[0].Action([]string{}, commands)
}

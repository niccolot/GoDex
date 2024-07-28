package main

import (
	"fmt"
	"bufio"
	"os"
	"os/exec"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}


func printPrompt() {

	fmt.Print("\nPokedex > ")
}

func printUnknown(text string) {

	fmt.Printf("'%s' command not found", text)
}

func commandHelp() error {
	fmt.Print("COMMAND MESSAGE HERE")
	return nil
}

func commandExit() error {
	// EXIT CLI
	return nil
}

func commandClear() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return nil
}

func cleanInput(text string) string {
	// removes trailing whitespaces and lowercases the command

	outCmd := strings.TrimSpace(text)
	outCmd = strings.ToLower(outCmd)
	return outCmd
}

func main() {

	cliCommandsTable := map[string]cliCommand{
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"exit": {
			name: "exit",
			description: "Quits the Pokedex CLI application and returns to terminal",
			callback: commandExit,
		},
		"clear": {
			name: "clear",
			description: "Clears the screen",
			callback: commandClear,
		},
	}

	reader := bufio.NewScanner(os.Stdin)
	printPrompt()
	for reader.Scan() {
		text := cleanInput(reader.Text())
		command, exists := cliCommandsTable[text]
		if exists {
			err := command.callback()
			if err != nil {
				fmt.Errorf("Failed to execute command '%s': %w", text, err)
			}
		} else {
			printUnknown(text)
		}
		printPrompt()
	}
}

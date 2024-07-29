package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"errors"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}


func printPrompt() {

	fmt.Print("\n\nPokedex > ")
}

func printUnknown(text string) {

	fmt.Printf("'%s' command not found", text)
}

func commandHelp() error {
	helpMessagePath := "help_message.txt"
	file, err := os.Open(helpMessagePath)

	if err != nil {
		return err
	} 

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	content := make([]byte, fileSize)
	_, err = file.Read(content)
	if err != nil {
		return err
	}
	fmt.Print(string(content))

	return nil
}

func commandExit() error {

	return nil
}

func commandClear() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return nil
}

func commandMap() error {
	locations := "https://pokeapi.co/api/v2/location-area/"
	res, err := http.Get(locations)
	if err != nil {
		return err
	}
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	
	defer res.Body.Close()

	if res.StatusCode > 299 {
		errorMsg := fmt.Sprintf("Response failsed with status code %d and\nbody: %s\n", res.StatusCode, body)
		return errors.New(errorMsg)
	}

	fmt.Print(len(body))

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
		"map": {
			name: "map",
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback: commandMap,
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
			if command.name == "exit" { return }
		} else {
			printUnknown(text)
		}
		printPrompt()
	}
}

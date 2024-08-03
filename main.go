package main

import (
	"fmt"
	"os"
	"bufio"
)

func main() {

	cliCommandsTable := map[string]CliCommand{
		"help": {
			Name: "help",
			Description: "Displays a help message",
			Callback: CommandHelp,
		},
		"exit": {
			Name: "exit",
			Description: "Quits the Pokedex CLI application and returns to terminal",
			Callback: CommandExit,
		},
		"clear": {
			Name: "clear",
			Description: "Clears the screen",
			Callback: CommandClear,
		},
		"map": {
			Name: "map",
			Description: "Displays the Names of 20 location areas in the Pokemon world",
			Callback: CommandMap,
		},
		"mapb": {
			Name: "umap",
			Description: "Displays the Names of the previous 20 location areas in the Pokemon world",
			Callback: CommandMapb,
		},
	}

	c := Config{
		LocationLimit: 10,
		LocationOffset: 0,
		PrevLocations: "",
	}
	c.CurrLocations = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=%d", 
									c.LocationOffset, 
									c.LocationLimit)
									
	c.NextLocations = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=%d", 
									c.LocationOffset + c.LocationLimit, 
									c.LocationLimit)

	reader := bufio.NewScanner(os.Stdin)
	PrintPrompt()
	for reader.Scan() {
		text := CleanInput(reader.Text())
		command, exists := cliCommandsTable[text]
		if exists {
			err := command.Callback(&c)
			if err != nil {
				fmt.Errorf("Failed to execute command '%s': %w", text, err)
			}
			if command.Name == "exit" { return }
		} else {
			PrintUnknown(text)
		}
		PrintPrompt()
	}
}

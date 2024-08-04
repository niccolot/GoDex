package main

import (
	"fmt"
	"strings"
)

func PrintUnknown(text string) {
	fmt.Printf("'%s' command not found", text)
}

func CleanInput(text string) string {
	// removes trailing whitespaces and lowercases the command
	outCmd := strings.TrimSpace(text)
	outCmd = strings.ToLower(outCmd)
	return outCmd
}

func getCliCommandsTable() map[string]CliCommand {
	table := map[string]CliCommand{
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
	
	return table
}

func getInitConfig() Config {
	config := Config{
		PrevLocations: "",
		NextLocations: "https://pokeapi.co/api/v2/location-area",
		History: make([]string, 10),
	}

	return config
}
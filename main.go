package main

import (
	"fmt"
	"github.com/peterh/liner"
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
		PrevLocations: "",
		NextLocations: "https://pokeapi.co/api/v2/location-area",
	}

	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	for {
		fmt.Println()
		input, err := line.Prompt("Pokedex > ")
		if err != nil {
			if err == liner.ErrPromptAborted {
				break
			}
			fmt.Println("Error reading line:", err)
			continue
		}

		text := CleanInput(input)
		command, exists := cliCommandsTable[text]
		if exists {
			err := command.Callback(&c)
			if err != nil {
				fmt.Println(fmt.Errorf("Failed to execute command '%s': %w", text, err).Error())
			}
			if command.Name == "exit" {
				break
			}
		} else {
			PrintUnknown(text)
		}

		line.AppendHistory(text)
		fmt.Println()
	}
	
}

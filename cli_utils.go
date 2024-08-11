package main

import (
	"fmt"
	"strings"
	"time"
	"math/rand"
	"github.com/niccolot/GoDex/internal/pokecache"
	"github.com/peterh/liner"
)


func PrintUnknown(text string) {
	fmt.Printf("'%s' command not found", text)
}

func ParseInput(text string) (command string, args []string) {
	/*
	* @param text (string): the whole input from terminal
	*
	* @return command, args (string, []string): the command name and optional arguments slice 
	*/

	// removes trailing whitespaces and lowercases the command
	trimmedText := strings.TrimSpace(text)
	lowercasedText := strings.ToLower(trimmedText)
	
	parts := strings.Split(lowercasedText, " ")

	command = parts[0]
	args = parts[1:]

	return command, args
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
		"history": {
			Name: "history",
			Description: "Displays the used commands",
			Callback: CommandHistory,
		},
		"explore": {
			Name: "explore",
			Description: "Displays the pokemons present in a certain area",
			Callback: CommandExplore,
		},
		"catch": {
			Name: "catch",
			Description: "Attempt to catch a pokemon",
			Callback: CommandCatch,
		},
		"inspect": {
			Name: "inspect",
			Description: "Displays information about a previously catch pokemon",
			Callback: CommandInspect,
		},
		"pokedex": {
			Name: "pokedex",
			Description: "Displays the pokemons contained in the pokedex",
			Callback: CommandPokedex,
		},
	}
	
	return table
}

func getInitConfig() *Config {
	minutesInCacheCommands := 5 * time.Minute
	minutesEscapedPokemon := 1 * time.Minute
	config := Config{
		PrevLocations: "",
		NextLocations: "https://pokeapi.co/api/v2/location-area",
		CurrLocation: "",
		ExplorableLocations: make([]string, 0),
		NearbyPokemons: make([]string, 0),
		History: make([]string, 0),
		PokeCache: *pokecache.NewCache[[]byte](minutesInCacheCommands),
		EscapedPokemons: *pokecache.NewCache[bool](minutesEscapedPokemon),
		Pokedex: make(map[string]PokeAPIPokemonInfo),
		Actions: map[string]CliCommand{
			"escape": {
				Name: "escape",
				Description: "Run away from a random encounter",
				Callback: CommandEscape,
			},
			"battle": {
				Name: "battle",
				Description: "Fight a pokemon found in a random encounter",
				Callback: CommandBattle,
			},
			"catch": {
				Name: "catch",
				Description: "Attempt to catch a pokemon",
				Callback: CommandCatch,
			},
			"exit": {
				Name: "exit",
				Description: "Quits the Pokedex CLI application and returns to terminal",
				Callback: CommandExit,
			},
		},
		EncounteredPokemon: "",
		BattleActions: map[string]CliCommand{
			"inspect": {
				Name: "inspect",
				Description: "Displays information about a previously catch pokemon",
				Callback: CommandInspect,
			},
			"pokedex": {
				Name: "pokedex",
				Description: "Displays the pokemons contained in the pokedex",
				Callback: CommandPokedex,
			},
			"choose": {
				Name: "choose",
				Description: "Choose a pokemon to fight with",
				Callback: CommandChoose,
			},
			"exit": {
				Name: "exit",
				Description: "Quits the Pokedex CLI application and returns to terminal",
				Callback: CommandExit,
			},
		},
	}

	return &config
}

func RandomEncounter(c *Config) error {
	if rand.Float64() > 0.0 {
		err := HandleRandomEncounter(c)
		if err != nil {
			return err
		}
	} 

	return nil
}

func HandleRandomEncounter(c *Config) error {
	encounteredPokemon := c.NearbyPokemons[rand.Intn(len(c.NearbyPokemons))]
	fmt.Printf("\nA wild %s appears!\n", encounteredPokemon)
	c.EncounteredPokemon = encounteredPokemon
	fmt.Println("Choose an action:")
	PrintActions(c)

	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	for {
		fmt.Println()
		input, err := line.Prompt("Pokedex/Encounter > ")
		if err != nil {
			return fmt.Errorf("failed reading line: %w", err)
		}

		go line.AppendHistory(input)
		c.History = append(c.History, input)

		input = input + " " + encounteredPokemon

		actionName, args := ParseInput(input)
		action, exists := c.Actions[actionName]
		if exists {
			err := action.Callback(c, args)
			if err != nil {
				
				fmt.Println(fmt.Errorf("failed to execute command '%s': %w", actionName, err).Error()) 
			}
			if action.Name == "escape" || action.Name == "exit"{
				break
			}
		} else {
			PrintUnknown(actionName)
		}
		fmt.Println()
	}
	
	return nil
}

func PrintActions(c *Config) {
	for action := range c.Actions {
		if action == "exit" {
			continue
		}
		fmt.Printf("- %s\n", action)
	}
}
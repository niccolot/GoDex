package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
	"github.com/peterh/liner"
)


func CommandHelp(c *Config, args []string) error {
	helpMessagePath := "assets/help_message.txt"
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

func CommandExit(c *Config, args []string) error {
	return nil
}

func CommandClear(c *Config, args []string) error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	
	return nil
}

func CommandMap(c *Config, args []string) error {
	locations := c.NextLocations

	body, err := GetData(c, locations)
	if err != nil {
		return err
	}

	next, prev, err := PrintLocations(c, body)
	c.NextLocations = next
	c.PrevLocations = prev

	return err
}

func CommandMapb(c *Config, args []string) error {
	locations := c.PrevLocations

	if locations == "" {
		return errors.New("no previous locations")
	}

	body, err := GetData(c, locations)
	if err != nil {
		return err
	}

	next, prev, err := PrintLocations(c, body)
	c.NextLocations = next
	c.PrevLocations = prev

	return err
}

func CommandHistory(c *Config, args []string) error {
	for _, entry := range c.History {
		fmt.Println(entry)
	}

	return nil
}

func CommandExplore(c *Config, args []string) error {
	if len(args) != 1 {
		return errors.New("command usage: explore <area-name>")
	}

	area := args[0]

	areaURL := areaURLAPI + area
	c.CurrLocation = areaURL

	if IsAreaNearby(c, area) {
		body, err := GetData(c, areaURL)
		if err != nil {
			return err
		}

		err = PrintPokemons(c, body)

		return err
	} else {
		fmt.Println("Only nearby areas can be explored!")
		fmt.Println("Use the 'map' and 'mapb' commands to move around in the world of pokemon")
		return nil
	}	
}

func CommandCatch(c *Config, args []string) error {
	if len(args) < 1 {
		return errors.New("command usage: catch <pokemon-name>")
	}

	pokemon := args[0]

	_, inPokedex := c.Pokedex[pokemon]
	if inPokedex {
		fmt.Printf("%s already in the Pokedex", pokemon)
		return nil
	}

	if c.CurrLocation == "" {
		fmt.Println("In order to catch some pokemon you first have to explore some areas!")
		fmt.Println("Use the command 'explore <area-name>' to find some pokemons")

		return nil
	}

	pokemonURL := pokemonURLAPI + pokemon
	pokemonStruct, err := GetPokemonStruct(c, pokemonURL)
	if err != nil {
		return err
	}

	if IsPokemonNearby(c, pokemon) {
		escaped, _ := c.EscapedPokemons.Get(pokemon)
		if  escaped {
			fmt.Printf("%s is still on the run! Try again in a while", pokemon)
			return nil
		}

		exp := pokemonStruct.BaseExperience

		fmt.Printf("Throwing a pokeball at %s...\n", pokemon)
		r := rand.Float64()

		// 340 listed as maximum base experience level
		if r > float64(exp)/340.0 {
			fmt.Printf("%s was caugth and added to the pokedex!", pokemon)
			c.Pokedex[pokemon] = pokemonStruct
		} else {
			go c.EscapedPokemons.Add(pokemon, true)
			fmt.Printf("%s escaped!", pokemon)
		}
	
	} else {
		fmt.Println("Only nearby pokemons can be caught!")
		fmt.Println("Use the command 'explore <area-name>' to list the pokemons near you")
	}
	return nil
}

func CommandInspect(c *Config, args []string) error {
	if len(args) != 1 {
		return errors.New("command usage: catch <pokemon-name>")
	}

	pokemon := args[0]
	pokemonStruct, inPokedex := c.Pokedex[pokemon]
	if !inPokedex {
		fmt.Printf("%s not found in pokedex, catch it in order to obtain some information on it", pokemon)
		return nil
	}
	PrintPokemonInfo(&pokemonStruct)

	return nil
}

func CommandPokedex(c *Config, args []string) error {
	if len(c.Pokedex) == 0 {
		fmt.Println("The pokedex is empty, go and catch some pokemons!")
		return nil
	}

	for key := range c.Pokedex {
		fmt.Printf("- %s\n", key)
	}

	return nil
}

func CommandEscape(c *Config, args []string) error {
	fmt.Println("Escaping...")

	return nil
}

func CommandBattle(c *Config, args []string) error {
	if len(c.Pokedex) == 0 {
		fmt.Println("You need to capture some pokemons in order to fight!")
		fmt.Println("Escaping...")
		return nil
	}

	if len(args) < 1 {
		return errors.New("command usage: battle <pokemon-name>")
	}

	pokemon := args[0]

	fmt.Printf("Choose a pokemon to fight with %s\n", pokemon)
	fmt.Printf("- Enter 'inspect %s' if you have already catch it to check its stats\n", pokemon)
	fmt.Println("- Enter 'pokedex' to check your pokedex")
	fmt.Println("- Enter 'inspect <pokemon-name>' to check the stats of one of your pokemons")
	fmt.Printf("- Enter 'choose <pokemon-name>' to start the battle with the chosen pokemon\n")

	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	for {
		fmt.Println()
		input, err := line.Prompt("Pokedex/Battle > ")
		if err != nil {
			return fmt.Errorf("failed reading line: %w", err)
		}

		line.AppendHistory(input)
		c.History = append(c.History, input)

		commandName, args := ParseInput(input)
		command, exists := c.BattleActions[commandName]
		if exists {
			err := command.Callback(c, args)
			if err != nil {
				fmt.Println(fmt.Errorf("failed to execute command '%s': %w", commandName, err).Error())
			}
			if command.Name == "exit" {
				break
			}
		} else {
			PrintUnknown(commandName)
		}
		fmt.Println()
	}
	
	return nil
}

func CommandChoose(c *Config, args []string) error {
	if len(args) != 1 {
		return errors.New("command usage: choose <pokemon-name>")
	}
	
	playerPokemonName := args[0]
	
	encounteredPokemonURL := pokemonURLAPI + c.EncounteredPokemon
	adversaryPokemon, err := GetPokemonStruct(c, encounteredPokemonURL)
	if err != nil {
		return err
	}

	playerPokemonURL := pokemonURLAPI + playerPokemonName
	playerPokemon, err := GetPokemonStruct(c, playerPokemonURL)
	if err != nil {
		return err
	}
	
	playerStats := GetPokemonStats(&playerPokemon)
	adversaryStats := GetPokemonStats(&adversaryPokemon)

	playerHp := playerStats.Hp
	adversaryHp := adversaryStats.Hp
	playerAttack := playerStats.Attack
	adversaryAttack := adversaryStats.Attack
	playerDefense := playerStats.Defense
	adversaryDefense := adversaryStats.Defense

	damageToAdversary := playerAttack - adversaryDefense
	if damageToAdversary < 0 {
		damageToAdversary = 0
	}

	damageToPlayer := adversaryAttack - playerDefense
	if damageToPlayer < 0 {
		damageToPlayer = 0
	}

	playerHp -= damageToPlayer
	adversaryHp -= damageToAdversary

	fmt.Printf("%s attacks %s for %d damage points\n", playerPokemonName, c.EncounteredPokemon, damageToAdversary)
	if adversaryHp < 0 {
		fmt.Printf("%s is stunned and got catched!\n", c.EncounteredPokemon)
		c.Pokedex[c.EncounteredPokemon] = adversaryPokemon

		return nil
	}

	fmt.Printf("%s attacks %s for %d damage points\n", c.EncounteredPokemon, playerPokemonName, damageToPlayer)
	if playerHp < 0 {
		fmt.Printf("%s is stunned and got removed from the pokedex!\n", playerPokemonName)
		delete(c.Pokedex, playerPokemonName)

		return nil
	}

	if playerHp > adversaryHp {
		fmt.Printf("%s won! %s escapes scared\n", playerPokemonName, c.EncounteredPokemon)
	} else {
		fmt.Printf("%s loses! Better escape while you can\n", playerPokemonName)
	}

	return nil
}

func CommandSave(c *Config, args []string) error {
	_, err := os.Stat("saves")
	if os.IsNotExist(err) {
		err := os.Mkdir("saves", os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create folder: %v", err)
		}
		fmt.Println("Created folder 'saves'")
	} else if err != nil {
		return fmt.Errorf("failed to check folder existence: %v", err)
	}

	currentTime := time.Now()

	// format the time as "day_month_year_hour-minute"
	folderName := currentTime.Format("02_01_2006_15-04")
	path := "saves/" + folderName
	filePath := path + ".json"

	err = SaveMapAsJSON(filePath, c.Pokedex)

	return err
}

func CommandLoad(c *Config, args []string) error {
	_, err := os.Stat("saves")
	if os.IsNotExist(err) {
		return fmt.Errorf("'saves' not found")
	}

	empty, err := IsFolderEmpty("saves")
	if err != nil {
		return err
	}
	if empty {
		fmt.Println("'saves' folder is empty")
		return nil
	} else {
		saves, err := GetFiles("saves")
		if err != nil {
			return err
		}

		fmt.Println("Enter one of the saved files: ")
		for _, save := range saves {
			fmt.Printf("- %s\n", save)
		}

		line := liner.NewLiner()
		defer line.Close()
		line.SetCtrlCAborts(true)
		for {
			fmt.Println()
			input, err := line.Prompt("Pokedex/Load > ")
			if err != nil {
				return fmt.Errorf("failed reading line: %w", err)
			}
	
			file, _ := ParseInput(input)
			if file == "exit" {
				break
			}
			exists := Contains(saves, file)
			if exists {
				filePath := "saves/" + file
				data, err := LoadMapFromJSON(filePath)
				if err != nil {
					return fmt.Errorf("failed to load %s file: %w", filePath, err)
				}
				c.Pokedex = data
				fmt.Println("Saved file loaded succesfully")
				break
			} else {
				fmt.Printf("Invalid choice\n")
			}
			fmt.Println()
		}
	}

	return nil
}
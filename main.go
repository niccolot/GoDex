package main

import (
	"fmt"
	"os"
	"bufio"
	"github.com/niccolot/GoDex/types"
	"github.com/niccolot/GoDex/commands"
	"github.com/niccolot/GoDex/utils"
)
/*
type cliCommand struct {
	Name        string
	Description string
	Callback    func(c *config) error
}

type config struct {
	locationLimit int
	locationOffset int
	prevLocations string
	nextLocations string
}

func printPrompt() {

	fmt.Print("\n\nPokedex > ")
}

func printUnknown(text string) {

	fmt.Printf("'%s' command not found", text)
}

func commandHelp(c *config) error {
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

func commandExit(c *config) error {

	return nil
}

func commandClear(c *config) error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return nil
}

func commandMap(c *config) error {
	locations := c.nextLocations
	c.prevLocations = locations
	c.locationOffset += c.locationLimit
	c.nextLocations = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=%d", 
									c.locationOffset, 
									c.locationLimit)
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
		errorMsg := fmt.Sprintf("Response failed with status code %d\n", res.StatusCode)
		return errors.New(errorMsg)
	}

	data := pokeapi.PokeAPIDataLocations{}
	errUnmarshal := json.Unmarshal(body, &data)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	fmt.Print(data.Results[0])

	return nil
}

func commandMapb(c *config) error {
	locations := c.prevLocations
	if locations == "" {
		return errors.New("No previous locations")
	}
	c.nextLocations = locations
	c.locationOffset -= c.locationLimit
	if c.locationOffset < 0 {
		c.locationOffset = 0
	}
	c.prevLocations = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=%d", 
									c.locationOffset, 
									c.locationLimit)
	
	err := printLocations(locations)

	return err
}

func printLocations(locations string) error {
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
		errorMsg := fmt.Sprintf("Response failed with status code %d\n", res.StatusCode)
		return errors.New(errorMsg)
	}

	data := pokeapi.PokeAPIDataLocations{}
	errUnmarshal := json.Unmarshal(body, &data)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	fmt.Print(data.Results[0])

	return nil
}

func cleanInput(text string) string {
	// removes trailing whitespaces and lowercases the command

	outCmd := strings.TrimSpace(text)
	outCmd = strings.ToLower(outCmd)
	return outCmd
}

*/

func main() {

	cliCommandsTable := map[string]types.CliCommand{
		"help": {
			Name: "help",
			Description: "Displays a help message",
			Callback: commands.CommandHelp,
		},
		"exit": {
			Name: "exit",
			Description: "Quits the Pokedex CLI application and returns to terminal",
			Callback: commands.CommandExit,
		},
		"clear": {
			Name: "clear",
			Description: "Clears the screen",
			Callback: commands.CommandClear,
		},
		"map": {
			Name: "map",
			Description: "Displays the Names of 20 location areas in the Pokemon world",
			Callback: commands.CommandMap,
		},
		"mapb": {
			Name: "umap",
			Description: "Displays the Names of the previous 20 location areas in the Pokemon world",
			Callback: commands.CommandMapb,
		},
	}

	c := types.Config{
		LocationLimit: 10,
		LocationOffset: 0,
		PrevLocations: "",
	}
	c.NextLocations = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=%d", 
									c.LocationOffset, 
									c.LocationLimit)

	reader := bufio.NewScanner(os.Stdin)
	utils.PrintPrompt()
	for reader.Scan() {
		text := utils.CleanInput(reader.Text())
		command, exists := cliCommandsTable[text]
		if exists {
			err := command.Callback(&c)
			if err != nil {
				fmt.Errorf("Failed to execute command '%s': %w", text, err)
			}
			if command.Name == "exit" { return }
		} else {
			utils.PrintUnknown(text)
		}
		utils.PrintPrompt()
	}
}

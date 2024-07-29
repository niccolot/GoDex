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
	"encoding/json"
	"github.com/niccolot/GoDex/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
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

func commandUMap(c *config) error {
	locations := c.prevLocations
	if locations == "" {
		return errors.New("No previous locations")
	}
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

	c := config{
		locationLimit: 10,
		locationOffset: 0,
		prevLocations: "",
	}
	c.nextLocations = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=%d", 
									c.locationOffset, 
									c.locationLimit)

	reader := bufio.NewScanner(os.Stdin)
	printPrompt()
	for reader.Scan() {
		text := cleanInput(reader.Text())
		command, exists := cliCommandsTable[text]
		if exists {
			err := command.callback(&c)
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

package main

import (
	"fmt"
	"errors"
	"os"
	"os/exec"
)


func CommandHelp(c *Config, args []string) error {
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

	body, err := GetLocationsBody(c, locations)
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

	body, err := GetLocationsBody(c, locations)
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

	areaURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", area)

	body, err := GetPokemons(c, areaURL)
	if err != nil {
		return err
	}

	err = PrintPokemons(body)

	return err
}
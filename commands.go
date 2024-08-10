package main

import (
	"fmt"
	"errors"
	"os"
	"os/exec"
	"net/http"
	"encoding/json"
	"io"
)

// ################## COMMANDS FUNCTIONS #################

func CommandHelp(c *Config) error {
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

func CommandExit(c *Config) error {

	return nil
}

func CommandClear(c *Config) error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	
	return nil
}

func CommandMap(c *Config) error {
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

func CommandMapb(c *Config) error {
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

func CommandHistory(c *Config) error {

	for _, entry := range c.History {
		fmt.Println(entry)
	}

	return nil
}

// ##################### UTILITY FUNCTIONS #################

func GetLocationsBody(c *Config, locations string) (body []byte, err error) {
	body, found := c.PokeCache.Get(locations)
	if found {

		return body, nil

	} else {
		res, err := http.Get(locations)
		if err != nil {
			return nil, err
		}
		
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		
		defer res.Body.Close()

		if res.StatusCode > 299 {
			errorMsg := fmt.Sprintf("Response failed with status code %d\n", res.StatusCode)
			return nil, errors.New(errorMsg)
		}

		go c.PokeCache.Add(locations, body)
		
		return body, nil
	}
}

func PrintLocations(c *Config, body []byte) (next string, prev string, err error) {
	data := PokeAPIDataLocations{}
	errUnmarshal := json.Unmarshal(body, &data)
	if errUnmarshal != nil {
		return c.NextLocations, c.PrevLocations, errUnmarshal
	}

	for _, entry := range data.Results {
		fmt.Println(entry.Name)
	}

	next = data.Next
	if data.Previous == nil {
		prev = ""
	} else {
		prev = *(data.Previous)
	}

	return next, prev, nil
}

/*
func PrintLocations(c *Config, locations string) (next string, prev string, err error) {
	res, err := http.Get(locations)
	if err != nil {
		return c.NextLocations, c.PrevLocations, err
	}
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return c.NextLocations, c.PrevLocations, err
	}
	
	defer res.Body.Close()

	if res.StatusCode > 299 {
		errorMsg := fmt.Sprintf("Response failed with status code %d\n", res.StatusCode)
		return c.NextLocations, c.PrevLocations, errors.New(errorMsg)
	}

	data := PokeAPIDataLocations{}
	errUnmarshal := json.Unmarshal(body, &data)
	if errUnmarshal != nil {
		return c.NextLocations, c.PrevLocations, errUnmarshal
	}

	for _, entry := range data.Results {
		fmt.Println(entry.Name)
	}

	next = data.Next
	if data.Previous == nil {
		prev = ""
	} else {
		prev = *(data.Previous)
	}

	return next, prev, nil
}
*/


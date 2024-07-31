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
	c.PrevLocations = locations
	c.LocationOffset += c.LocationLimit
	c.NextLocations = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=%d", 
									c.LocationOffset, 
									c.LocationLimit)
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

	data := PokeAPIDataLocations{}
	errUnmarshal := json.Unmarshal(body, &data)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	fmt.Print(data.Results[0])

	return nil
}

func CommandMapb(c *Config) error {
	locations := c.PrevLocations
	if locations == "" {
		return errors.New("No previous locations")
	}
	c.NextLocations = locations
	c.LocationOffset -= c.LocationLimit
	if c.LocationOffset < 0 {
		c.LocationOffset = 0
	}
	c.PrevLocations = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=%d", 
									c.LocationOffset, 
									c.LocationLimit)
	
	err := PrintLocations(locations)

	return err
}

func PrintLocations(locations string) error {
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

	data := PokeAPIDataLocations{}
	errUnmarshal := json.Unmarshal(body, &data)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	fmt.Print(data.Results[0])

	return nil
}
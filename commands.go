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
	fmt.Printf("\n\nMAP Offset before: %d\n", c.LocationOffset)
	fmt.Printf("MAP prev locations before: %s\n", c.PrevLocations)
	fmt.Printf("MAP curr locations before: %s\n", c.CurrLocations)
	fmt.Printf("MAP  next locations before: %s\n", c.NextLocations)
	locations := c.NextLocations
	c.PrevLocations = c.CurrLocations
	c.CurrLocations = c.NextLocations
	c.LocationOffset += c.LocationLimit
	c.NextLocations = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=%d", 
									c.LocationOffset, 
									c.LocationLimit)
	fmt.Printf("\n\nMAP Offset after: %d\n", c.LocationOffset)
	fmt.Printf("MAP prev locations after: %s\n", c.PrevLocations)
	fmt.Printf("MAP curr locations after: %s\n", c.CurrLocations)
	fmt.Printf("MAP  next locations after: %s\n\n", c.NextLocations)
	err := PrintLocations(locations)

	return err
}

func CommandMapb(c *Config) error {
	fmt.Printf("\n\nMAPB Offset before: %d\n", c.LocationOffset)
	fmt.Printf("MAPB prev locations before: %s\n", c.PrevLocations)
	fmt.Printf("MAPB curr locations before: %s\n", c.CurrLocations)
	fmt.Printf("MAPB  next locations before: %s\n", c.NextLocations)
	locations := c.PrevLocations
	if locations == "" {
		return errors.New("no previous locations")
	}
	c.NextLocations = c.CurrLocations
	c.CurrLocations = c.PrevLocations
	c.LocationOffset -= c.LocationLimit
	if c.LocationOffset < 0 {
		c.LocationOffset = 0
		c.PrevLocations = ""
		c.CurrLocations = ""
	} else {
		c.PrevLocations = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=%d", 
										c.LocationOffset, 
										c.LocationLimit)
	}
	fmt.Printf("\n\nMAPB Offset after: %d\n", c.LocationOffset)
	fmt.Printf("MAPB prev locations after: %s\n", c.PrevLocations)
	fmt.Printf("MAPB curr locations after: %s\n", c.CurrLocations)
	fmt.Printf("MAPB  next locations after: %s\n\n", c.NextLocations)
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
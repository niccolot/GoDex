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
	next, prev, err := PrintLocations(c, locations)
	c.NextLocations = next
	c.PrevLocations = prev

	return err
}

func CommandMapb(c *Config) error {
	locations := c.PrevLocations
	if locations == "" {
		return errors.New("no previous locations")
	}
	next, prev, err := PrintLocations(c, locations)
	c.NextLocations = next
	c.PrevLocations = prev

	return err
}

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

	fmt.Print(data.Results[0])
	next = data.Next
	if data.Previous == nil {
		prev = ""
	} else {
		prev = *(data.Previous)
	}

	return next, prev, nil
}
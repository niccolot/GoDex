package main

import (
	"fmt"
	"errors"
	"net/http"
	"encoding/json"
	"io"
)


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

func GetPokemons(c *Config, area string) (body []byte, err error) {
	body, found := c.PokeCache.Get(area)
	if found {

		return body, nil

	} else {
		res, err := http.Get(area)
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

		go c.PokeCache.Add(area, body)
		
		return body, nil
	}
}

func PrintPokemons(body []byte) error {
	data := PokeAPIAreaInfo{}
	errUnmarshal := json.Unmarshal(body, &data)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	encounters := data.PokemonEncounters

	for _, encounter := range encounters {
		fmt.Println("- ", encounter.Pokemon.Name)
	}

	return nil
}
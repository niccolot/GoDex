package main

import (
	"fmt"
	"errors"
	"net/http"
	"encoding/json"
	"io"
)


func GetData(c *Config, url string) (body []byte, err error) {
	/*
	* given an url representing some data (e.g. areas or pokemons) 
	* that could be either in the cache or to be fetched from the
	* webserver, it returns the body of the data and an error object
	*
	* @param c (*Config): the pointer to the commands configuration struct
	* @param url (string): the url belonging to the resources in demand
	* @return body ([]byte): the data as a raw byte slice, nil if an error occurred
	* @return err (error): eventual error occurred during the fetching of the data,
		nil if everything went well
	*/

	body, found := c.PokeCache.Get(url)
	if found {
		return body, nil

	} else {
		body, err := GetBodyFromHTTP(url)
		if err != nil {
			return nil, err
		}

		go c.PokeCache.Add(url, body)
		
		return body, nil
	}
}

func GetBodyFromHTTP(url string) (body []byte, err error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	
	defer res.Body.Close()

	if res.StatusCode > 299 {
		errorMsg := fmt.Sprintf("Response failed with status code %d\n", res.StatusCode)
		return nil, errors.New(errorMsg)
	}

	return body, err
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

func GetPokemonStruct(pokemon string) (pokemonStruct PokeAPIPokemonInfo, err error) {
	body, err := GetBodyFromHTTP(pokemon)
	if err != nil {
		return PokeAPIPokemonInfo{}, err
	}
	pokemonStruct = PokeAPIPokemonInfo{}
	errUnmarshal := json.Unmarshal(body, &pokemonStruct)
	if errUnmarshal != nil {
		return PokeAPIPokemonInfo{}, errUnmarshal
	}

	return pokemonStruct, nil
}

func PrintPokemonInfo(pokemon *PokeAPIPokemonInfo) {
	name := pokemon.Name
	height := pokemon.Height
	weight := pokemon.Weight
	stats := pokemon.Stats
	types := pokemon.Types

	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Height: %d\n", height)
	fmt.Printf("Weight: %d\n", weight)
	
	fmt.Println("Stats:")
	for _, s := range stats {
		fmt.Printf(" -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	
	fmt.Println("Types:")
	for _, t := range types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}
}
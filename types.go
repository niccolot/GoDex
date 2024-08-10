package main

import (
	"github.com/niccolot/GoDex/internal/pokecache"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(c *Config, args []string) error
}

type Config struct {
	PrevLocations string
	NextLocations string
	History []string
	PokeCache pokecache.Cache
	Pokedex map[string]PokeAPIPokemonInfo
}
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
	CurrLocation string
	NearbyPokemons []string
	ExplorableLocations []string
	History []string
	PokeCache pokecache.Cache[[]byte]
	EscapedPokemons pokecache.Cache[bool]
	Pokedex map[string]PokeAPIPokemonInfo
	Actions map[string]CliCommand
	EncounteredPokemon string
	BattleActions map[string]CliCommand
}
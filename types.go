package main

type CliCommand struct {
	Name        string
	Description string
	Callback    func(c *Config) error
}

type Config struct {
	LocationLimit int
	LocationOffset int
	PrevLocations string
	NextLocations string
	CurrLocations string
}
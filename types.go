package main

type CliCommand struct {
	Name        string
	Description string
	Callback    func(c *Config) error
}

type Config struct {
	PrevLocations string
	NextLocations string
	History []string
}
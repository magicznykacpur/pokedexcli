package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/magicznykacpur/pokedexcli/internal/pokeapi"
	"github.com/magicznykacpur/pokedexcli/internal/pokecache"
)

type CliCommand struct {
	name        string
	description string
	Callback    func(c *Config) error
}

type Config struct {
	Next     string
	Previous string
}

func GetSupportedCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 location areas in Pokemon world",
			Callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 location areas in Pokemon world",
			Callback:    commandMapB,
		},
	}
}

func commandExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(c *Config) error {
	help := "Welcome to the Pokedex!\nUsage:\n\n"
	supportedCommands := GetSupportedCommands()

	for _, command := range supportedCommands {
		help += fmt.Sprintf("%s: %s\n", command.name, command.description)
	}

	fmt.Println(help)

	return nil
}

var cache = pokecache.NewCache(time.Second * 69)

func useCachedResponse(cachedBytes []byte, c *Config) error {
	var locationArea pokeapi.LocationArea
	if err := json.Unmarshal(cachedBytes, &locationArea); err != nil {
		return fmt.Errorf("couldn't unmarshal location areas: %v", err)
	}

	for _, location := range locationArea.Results {
		fmt.Println(location.Name)
	}

	c.Previous = locationArea.Previous
	c.Next = locationArea.Next

	return nil
}

func commandMap(c *Config) error {
	var locationsAreaUrl string
	if c.Next != "" {
		locationsAreaUrl = c.Next
	} else {
		locationsAreaUrl = pokeapi.BaseUrl + "location-area/?offset=0&limit=20"
	}

	cachedBytes, ok := cache.Get(locationsAreaUrl)
	if ok {
		return useCachedResponse(cachedBytes, c)
	}

	bytes, err := pokeapi.GetLocationAreaBytes(locationsAreaUrl)
	if err != nil {
		return fmt.Errorf("could get location area bytes: %v", err)
	}

	locationArea, err := pokeapi.UnmarshalLocationArea(bytes)
	if err != nil {
		return err
	}

	cache.Add(locationsAreaUrl, bytes)

	c.Previous = locationArea.Previous
	c.Next = locationArea.Next

	for _, location := range locationArea.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapB(c *Config) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		c.Next = ""
		return nil
	}

	cachedBytes, ok := cache.Get(c.Previous)
	if ok {
		return useCachedResponse(cachedBytes, c)
	}

	bytes, err := pokeapi.GetLocationAreaBytes(c.Previous)
	if err != nil {
		return fmt.Errorf("couldn't get location area bytes: %v", err)
	}

	locationArea, err := pokeapi.UnmarshalLocationArea(bytes)
	if err != nil {
		return err
	}

	cache.Add(locationArea.Previous, bytes)
	c.Previous = locationArea.Previous
	c.Next = locationArea.Next

	for _, location := range locationArea.Results {
		fmt.Println(location.Name)
	}

	return nil
}

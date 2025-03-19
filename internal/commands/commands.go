package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

const baseUrl string = "https://pokeapi.co/api/v2/"

type locationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

func commandMap(c *Config) error {
	var locationsAreaUrl string
	if c.Next == "" {
		locationsAreaUrl = baseUrl + "location-area/"
	} else {
		locationsAreaUrl = c.Next
	}

	res, err := http.Get(locationsAreaUrl)
	if err != nil {
		return fmt.Errorf("couldn't get location areas: %v", err)
	}
	defer res.Body.Close()

	var locationArea locationArea
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&locationArea); err != nil {
		return fmt.Errorf("couldn't decode location areas: %v", err)
	}

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
		return nil
	}

	res, err := http.Get(c.Previous)
	if err != nil {
		return fmt.Errorf("couldn't get location areas: %v", err)
	}
	defer res.Body.Close()

	var locationArea locationArea
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&locationArea); err != nil {
		return fmt.Errorf("couldn't decode location area: %v", err)
	}

	c.Previous = locationArea.Previous
	c.Next = locationArea.Next

	for _, location := range locationArea.Results {
		fmt.Println(location.Name)
	}

	return nil
}

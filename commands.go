package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

type config struct {
	Next     string
	Previous string
}

func getSupportedCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 location areas in Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 location areas in Pokemon world",
			callback:    commandMapB,
		},
	}
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(c *config) error {
	help := "Welcome to the Pokedex!\nUsage:\n\n"
	supportedCommands := getSupportedCommands()

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

func commandMap(c *config) error {
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

func commandMapB(c *config) error {
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

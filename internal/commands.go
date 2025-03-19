package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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

var cache = NewCache(time.Second * 69)

func useCachedResponse(cachedBytes []byte, c *Config) error {
	var locationArea locationArea
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

func getLocationAreaBytes(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't get location areas: %v", err)
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't read bytes of response: %v", err)
	}

	return bytes, nil
}

func commandMap(c *Config) error {
	var locationsAreaUrl string
	if c.Next != "" {
		locationsAreaUrl = c.Next
	} else {
		locationsAreaUrl = baseUrl + "location-area/?offset=0&limit=20"
	}

	cachedBytes, ok := cache.Get(locationsAreaUrl)
	if ok {
		return useCachedResponse(cachedBytes, c)
	}

	bytes, err := getLocationAreaBytes(locationsAreaUrl)
	if err != nil {
		return fmt.Errorf("could get location area bytes: %v", err)
	}

	var locationArea locationArea
	if err := json.Unmarshal(bytes, &locationArea); err != nil {
		return fmt.Errorf("couldn't unmarshal location areas: %v", err)
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

	bytes, err := getLocationAreaBytes(c.Previous)
	if err != nil {
		return fmt.Errorf("couldn't get location area bytes: %v", err)
	}

	var locationArea locationArea
	if err := json.Unmarshal(bytes, &locationArea); err != nil {
		return fmt.Errorf("couldn't unmarshal location areas: %v", err)
	}

	cache.Add(locationArea.Previous, bytes)
	c.Previous = locationArea.Previous
	c.Next = locationArea.Next

	for _, location := range locationArea.Results {
		fmt.Println(location.Name)
	}

	return nil
}

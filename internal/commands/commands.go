package commands

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/magicznykacpur/pokedexcli/internal/decoding"
	"github.com/magicznykacpur/pokedexcli/internal/pokeapi"
	"github.com/magicznykacpur/pokedexcli/internal/pokecache"
	"github.com/magicznykacpur/pokedexcli/internal/pokedex"
)

type CliCommand struct {
	name        string
	description string
	Callback    func(c *Config, args ...string) error
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
		"explore": {
			name:        "explore",
			description: "Displays pokemons that can be found in the given area",
			Callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a pokemon",
			Callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Shows detailed info about a caught pokemon",
			Callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows all the caught pokemon",
			Callback:    commandPokedex,
		},
	}
}

func commandExit(c *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(c *Config, args ...string) error {
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
	locationArea, err := decoding.UnmarshalLocationArea(cachedBytes)
	if err != nil {
		return err
	}
	
	for _, location := range locationArea.Results {
		fmt.Println(location.Name)
	}

	c.Previous = locationArea.Previous
	c.Next = locationArea.Next

	return nil
}

func commandMap(c *Config, args ...string) error {
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
		return err
	}

	locationArea, err := decoding.UnmarshalLocationArea(bytes)
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

func commandMapB(c *Config, args ...string) error {
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
		return err
	}

	locationArea, err := decoding.UnmarshalLocationArea(bytes)
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

func commandExplore(_ *Config, args ...string) error {
	location := args[1]
	
	cachedBytes, ok := cache.Get(location)
	if ok {
		locationAreaByLocation, err := decoding.UnmarshalLocationAreaByLocation(cachedBytes)
		if err != nil {
			return err
		}
		
		names := ""
		for _, encounter := range locationAreaByLocation.PokemonEncounters {
			names += encounter.Pokemon.Name + " "
		}

		fmt.Println(names)
		return nil
	}

	bytes, err := pokeapi.GetLocationAreaByLocationBytes(location)
	if err != nil {
		return err
	}

	cache.Add(location, bytes)
	
	locationAreaByLocation, err := decoding.UnmarshalLocationAreaByLocation(bytes)
	if err != nil {
		return err
	}
	
	fmt.Printf("Exploring %s...\n", location)

	names := ""
	for _, encounter := range locationAreaByLocation.PokemonEncounters {
		names += fmt.Sprintf(" - %s\n", encounter.Pokemon.Name)
	}

	if len(names) > 0 {
		fmt.Println("Found pokemons:")
		fmt.Println(names)
	} else {
		fmt.Println("Nothing found...")
	}
	
	return nil
}

var usersPokedex = pokedex.NewPokedex()

func commandCatch(c *Config, args ...string) error {
	name := args[1]
	_, ok := usersPokedex.Get(name)
	if ok {
		return fmt.Errorf("you've already caught %s", name)
	}

	bytes, err := pokeapi.GetPokemonByName(name)
	if err != nil {
		return err
	}

	pokemon, err := decoding.UnmarshalPokemon(bytes)
	if err != nil {
		return err
	}

	baseExp := pokemon.BaseExperience
	randomInt := rand.IntN(baseExp)

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	if randomInt > baseExp / 2 {
		fmt.Printf("%s was caught!\n", name)
		fmt.Println("You may now inspect it with the inspect command.")
		usersPokedex.Catch(pokemon)
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
	
	return nil
}

func commandInspect(c *Config, args ...string) error {
	name := args[1]
	pokemon, ok := usersPokedex.Get(name)
	if !ok {
		return fmt.Errorf("you have not caught this pokemon")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Abilites:")
	for _, ability := range pokemon.Abilities {
		fmt.Printf("-%s\n", ability.Ability.Name)
	}

	fmt.Println("Types:")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("-%s\n", pokemonType.Type.Name)
	}

	return nil
}

func commandPokedex(c *Config, args ...string) error {
	if usersPokedex.IsEmpty() {
		return fmt.Errorf("your pokedex is empty")
	}

	fmt.Println("Your Pokedex:")
	for key, _ := range usersPokedex.GetCaughtPokemons() {
		fmt.Printf("- %s\n", key)
	}

	return nil
}
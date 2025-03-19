package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/magicznykacpur/pokedexcli/internal/commands"
	"github.com/magicznykacpur/pokedexcli/internal/pokecache"
)

func cleanInput(text string) []string {
	output := strings.TrimSpace(text)
	return strings.Split(strings.ToLower(output), " ")
}

func main() {
	supportedCommands := commands.GetSupportedCommands()
	scanner := bufio.NewScanner(os.Stdin)
	config := commands.Config{Previous: "", Next: ""}
	pokecache.NewCache(time.Millisecond * 1000)

	for {
		fmt.Printf("Pokedex > ")

		scanner.Scan()
		input := cleanInput(scanner.Text())

		command, ok := supportedCommands[input[0]]
		if ok {
			err := command.Callback(&config)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

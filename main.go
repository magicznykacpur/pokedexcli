package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/magicznykacpur/pokedexcli/internal"
)

func cleanInput(text string) []string {
	output := strings.TrimSpace(text)
	return strings.Split(strings.ToLower(output), " ")
}

func main() {
	supportedCommands := internal.GetSupportedCommands()
	scanner := bufio.NewScanner(os.Stdin)
	config := internal.Config{Previous: "", Next: ""}

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

package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	output := strings.TrimSpace(text)
	return strings.Split(output, " ")
}

func main() {
	fmt.Println("Hello, World!")
}
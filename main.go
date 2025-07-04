package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(r cmdRegistry) error
}

type cmdRegistry = map[string]cliCommand

func cleanInput(text string) []string {
	words := strings.Split(text, " ")
	result := make([]string, 0, len(words))
	for _, word := range words {
		trimmed := strings.Trim(word, " ")
		if trimmed == "" {
			continue
		}

		lowercased := strings.ToLower(trimmed)

		result = append(result, lowercased)
	}
	return result
}

func commandExit(r cmdRegistry) error {
	fmt.Println("\nClosing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(r cmdRegistry) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Printf("Usage: \n\n")
	for _, cmd := range r {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	registry := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		inputLine := scanner.Text()
		words := cleanInput(inputLine)
		command := words[0]

		if cmd, ok := registry[command]; ok {
			cmd.callback(registry)
		} else {
			fmt.Println("Unknown command")
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"gautam1168/pokedexcli/internal/pokeapi"
	"gautam1168/pokedexcli/internal/pokecache"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*state) error
}

type state struct {
	page        pokeapi.PokeLocationPage
	cmdRegistry map[string]cliCommand
	cache       *pokecache.Cache
}

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

func commandExit(s *state) error {
	fmt.Println("\nClosing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(s *state) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Printf("Usage: \n\n")
	for _, cmd := range s.cmdRegistry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(s *state) error {
	locationData, err := pokeapi.GetPokeLocations(&s.page, s.cache)
	locationData.Offset += 20
	if err != nil {
		return err
	} else {
		s.page = locationData
	}

	fmt.Println("")
	for _, location := range locationData.Locations {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapBack(s *state) error {
	if s.page.Offset == 0 {
		fmt.Println("\nyou're on the first page")
	} else {
		s.page.Offset -= 20
		locationData, err := pokeapi.GetPokeLocations(&s.page, s.cache)
		if err != nil {
			return err
		} else {
			s.page = locationData
		}

		fmt.Println("")
		for _, location := range locationData.Locations {
			fmt.Println(location.Name)
		}

		return nil
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
		"map": {
			name:        "map",
			description: "Browse the map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Rewind the map",
			callback:    commandMapBack,
		},
	}

	s := state{
		cmdRegistry: registry,
		page: pokeapi.PokeLocationPage{
			Offset: 0,
		},
		cache: pokecache.NewCache(5 * time.Second),
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		inputLine := scanner.Text()
		words := cleanInput(inputLine)
		command := words[0]

		if cmd, ok := registry[command]; ok {
			err := cmd.callback(&s)
			if err != nil {
				fmt.Println("\n", err.Error())
			}
		} else {
			fmt.Println("\nUnknown command")
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"gautam1168/pokedexcli/internal/pokeapi"
	"gautam1168/pokedexcli/internal/pokecache"
	"math/rand"
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
	pokemons    map[string]bool
	cache       *pokecache.Cache
	args        []string
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

func commandExplore(s *state) error {
	args := s.args
	if len(args) != 1 {
		return fmt.Errorf("expected exactly 1 argument for explore command but obtained: %v", len(args))
	}

	locations, err := pokeapi.GetPokemonInLocation(args[0], s.cache)
	if err != nil {
		return err
	}

	fmt.Printf("\nExploring %s...\n", args[0])
	fmt.Println("Found Pokemon:")
	for _, loc := range locations {
		fmt.Printf(" - %s\n", loc.Name)
	}

	return nil
}

func commandCatch(s *state) error {
	args := s.args
	if len(args) != 1 {
		return fmt.Errorf("expected exactly 1 argument for catch command but obtained: %v", len(args))
	}
	if pokemonDetails, err := pokeapi.GetPokemonDetails(args[0], s.cache); err != nil {
		return err
	} else {
		baseExperience := pokemonDetails.BaseExperience
		fmt.Printf("Throwing a Pokeball at %s...\n", args[0])
		if baseExperience == 0 {
			s.pokemons[args[0]] = true
			fmt.Printf("%s was caught!\n", args[0])
		} else {
			inverseBaseExp := (float64(baseExperience) / 400)
			if inverseBaseExp >= 1 {
				inverseBaseExp = 0.95
			}
			randomVal := float64(rand.Intn(100)) / 100
			// fmt.Printf("base: %v, experince: %v, chance: %v\n", baseExperience, inverseBaseExp, randomVal)
			if randomVal >= inverseBaseExp {
				s.pokemons[args[0]] = true
				fmt.Printf("%s was caught!\n", args[0])
			} else {
				fmt.Printf("%s escaped!\n", args[0])
			}
		}
	}
	return nil
}

func commandInspect(s *state) error {
	args := s.args
	if len(args) != 1 {
		return fmt.Errorf("expected exactly 1 argument for inspect command but obtained: %v", len(args))
	}

	if _, ok := s.pokemons[args[0]]; ok {
		if pokemonDetails, err := pokeapi.GetPokemonDetails(args[0], s.cache); err == nil {
			fmt.Println("Name: ", args[0])
			fmt.Println("Height: ", pokemonDetails.Height)
			fmt.Println("Weight: ", pokemonDetails.Weight)
			fmt.Println("Stats: ")
			for _, stat := range pokemonDetails.Stats {
				fmt.Printf(" -%s:%v\n", stat.Stat.Name, stat.BaseStat)
			}

			fmt.Println("Types:")
			for _, poktype := range pokemonDetails.Types {
				fmt.Printf(" - %s\n", poktype.Type.Name)
			}
		} else {
			return err
		}
	} else {
		fmt.Println("you have not caught that pokemon")
	}
	return nil
}

func commandPokedex(s *state) error {
	fmt.Println("Your Pokedex: ")
	for key := range s.pokemons {
		fmt.Printf(" - %s\n", key)
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
		"explore": {
			name:        "explore",
			description: "See pokemons in location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon by name",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect the stats of a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List pokemons you have caught",
			callback:    commandPokedex,
		},
	}

	s := state{
		cmdRegistry: registry,
		pokemons:    make(map[string]bool),
		page: pokeapi.PokeLocationPage{
			Offset: 0,
		},
		cache: pokecache.NewCache(5 * time.Second),
		args:  []string{},
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		inputLine := scanner.Text()
		words := cleanInput(inputLine)
		command := words[0]

		if cmd, ok := registry[command]; ok {
			if len(words) > 1 {
				s.args = words[1:]
			}

			err := cmd.callback(&s)
			if err != nil {
				fmt.Println("\n", err.Error())
			}
		} else {
			fmt.Println("\nUnknown command")
		}
	}
}

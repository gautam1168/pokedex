package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		inputLine := scanner.Text()
		words := cleanInput(inputLine)
		if len(words) > 0 {
			fmt.Printf("Your command was: %s\n", words[0])
		}
	}
}

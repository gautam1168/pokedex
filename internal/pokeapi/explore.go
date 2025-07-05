package pokeapi

import (
	"gautam1168/pokedexcli/internal/pokecache"
)

type Pokemon struct {
	Name string `json:"name"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type LocationDetails struct {
	Id         int                `json:"id"`
	Name       string             `json:"name"`
	Encounters []PokemonEncounter `json:"pokemon_encounters"`
}

func GetPokemonInLocation(location string, cache *pokecache.Cache) ([]Pokemon, error) {
	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	fullUrl := baseUrl + location

	result := []Pokemon{}

	if parsedLocation, err := GetDataAndParse[LocationDetails](fullUrl, cache); err != nil {
		return result, err
	} else {
		for _, poke := range parsedLocation.Encounters {
			result = append(result, poke.Pokemon)
		}
		return result, nil
	}
}

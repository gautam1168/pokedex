package pokeapi

import (
	"encoding/json"
	"gautam1168/pokedexcli/internal/pokecache"
	"io"
	"net/http"
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

	parsedLocation := LocationDetails{}
	if cachedBytes, ok := cache.Get(fullUrl); ok {
		if err := json.Unmarshal(cachedBytes, &parsedLocation); err != nil {
			return result, nil
		}
	} else {
		request, err := http.NewRequest("GET", fullUrl, nil)
		if err != nil {
			return result, err
		}

		request.Header.Set("Content-Type", "application/json")
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return result, err
		}
		defer response.Body.Close()
		if networkBytes, err := io.ReadAll(response.Body); err == nil {
			if err := json.Unmarshal(networkBytes, &parsedLocation); err != nil {
				return result, nil
			}
			cache.Add(fullUrl, networkBytes)
		} else {
			return result, err
		}
	}

	for _, poke := range parsedLocation.Encounters {
		result = append(result, poke.Pokemon)
	}

	return result, nil
}

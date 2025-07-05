package pokeapi

import "gautam1168/pokedexcli/internal/pokecache"

type stat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type pokemontype struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type PokemonDetail struct {
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []stat        `json:"stats"`
	Types          []pokemontype `json:"types"`
}

func GetPokemonDetails(name string, cache *pokecache.Cache) (PokemonDetail, error) {
	baseUrl := "https://pokeapi.co/api/v2/pokemon/"
	fullUrl := baseUrl + name

	pokemonDetails, err := GetDataAndParse[PokemonDetail](fullUrl, cache)
	return pokemonDetails, err
}

package pokeapi

import "gautam1168/pokedexcli/internal/pokecache"

type PokemonDetail struct {
	BaseExperience int `json:"base_experience"`
	Height         int `json:"height"`
	Weight         int `json:"weight"`
}

func GetPokemonDetails(name string, cache *pokecache.Cache) (PokemonDetail, error) {
	baseUrl := "https://pokeapi.co/api/v2/pokemon/"
	fullUrl := baseUrl + name

	pokemonDetails, err := GetDataAndParse[PokemonDetail](fullUrl, cache)
	return pokemonDetails, err
}

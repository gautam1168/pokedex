package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PokeLocation struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokeLocationData struct {
	Count     int            `json:"count"`
	Next      string         `json:"next"`
	Prev      string         `json:"previous"`
	Locations []PokeLocation `json:"results"`
}

type PokeLocationPage struct {
	Locations []PokeLocation
	Offset    int
}

func GetPokeLocations(page *PokeLocationPage) (PokeLocationPage, error) {
	if page == nil {
		return PokeLocationPage{}, fmt.Errorf("page cannot be nil")
	}

	result := PokeLocationPage{
		Offset: page.Offset,
	}

	baseUrl := "https://pokeapi.co/api/v2/location-area"
	fullUrl := fmt.Sprintf("%s?offset=%v&limit=%v", baseUrl, result.Offset, 20)
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

	decoder := json.NewDecoder(response.Body)
	apiData := PokeLocationData{}
	if err := decoder.Decode(&apiData); err != nil {
		return result, err
	}

	result.Locations = apiData.Locations

	return result, nil
}

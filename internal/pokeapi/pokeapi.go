package pokeapi

import (
	"fmt"
	"gautam1168/pokedexcli/internal/pokecache"
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

func GetPageUrl(page *PokeLocationPage) (string, error) {
	if page == nil {
		return "", fmt.Errorf("page cannot be nil")
	}

	baseUrl := "https://pokeapi.co/api/v2/location-area"
	fullUrl := fmt.Sprintf("%s?offset=%v&limit=%v", baseUrl, page.Offset, 20)
	return fullUrl, nil
}

func GetPokeLocations(page *PokeLocationPage, cache *pokecache.Cache) (PokeLocationPage, error) {
	if page == nil {
		return PokeLocationPage{}, fmt.Errorf("page cannot be nil")
	}

	if cache == nil {
		return PokeLocationPage{}, fmt.Errorf("pokecache cannot be nil")
	}

	result := PokeLocationPage{
		Offset: page.Offset,
	}

	fullUrl, err := GetPageUrl(page)
	if err != nil {
		return result, err
	}

	if apiData, err := GetDataAndParse[PokeLocationData](fullUrl, cache); err != nil {
		return result, err
	} else {
		result.Locations = apiData.Locations
	}

	return result, nil
}

// func GetDataAndParse[T any](fullUrl string, cache *pokecache.Cache) (T, error) {
// 	var apiData T
// 	if cachedBytes, ok := cache.Get(fullUrl); ok {
// 		if err := json.Unmarshal(cachedBytes, &apiData); err != nil {
// 			return apiData, err
// 		}
// 	} else {
// 		request, err := http.NewRequest("GET", fullUrl, nil)
// 		if err != nil {
// 			return apiData, err
// 		}

// 		request.Header.Set("Content-Type", "application/json")
// 		response, err := http.DefaultClient.Do(request)
// 		if err != nil {
// 			return apiData, err
// 		}
// 		defer response.Body.Close()

// 		networkBytes, err := io.ReadAll(response.Body)
// 		if err != nil {
// 			return apiData, err
// 		} else {
// 			cache.Add(fullUrl, networkBytes)
// 		}

// 		if err := json.Unmarshal(networkBytes, &apiData); err != nil {
// 			return apiData, err
// 		}
// 	}
// 	return apiData, nil
// }

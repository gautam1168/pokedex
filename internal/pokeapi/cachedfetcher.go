package pokeapi

import (
	"encoding/json"
	"gautam1168/pokedexcli/internal/pokecache"
	"io"
	"net/http"
)

func GetDataAndParse[T any](fullUrl string, cache *pokecache.Cache) (T, error) {
	var apiData T
	if cachedBytes, ok := cache.Get(fullUrl); ok {
		if err := json.Unmarshal(cachedBytes, &apiData); err != nil {
			return apiData, err
		}
	} else {
		request, err := http.NewRequest("GET", fullUrl, nil)
		if err != nil {
			return apiData, err
		}

		request.Header.Set("Content-Type", "application/json")
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return apiData, err
		}
		defer response.Body.Close()

		networkBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return apiData, err
		} else {
			cache.Add(fullUrl, networkBytes)
		}

		if err := json.Unmarshal(networkBytes, &apiData); err != nil {
			return apiData, err
		}
	}
	return apiData, nil
}

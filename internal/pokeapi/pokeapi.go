package pokeapi

import (
	"fmt"
	"io"
	"net/http"
)

const BaseUrl string = "https://pokeapi.co/api/v2/"

func GetLocationAreaBytes(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't get location area: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("couldn't get location area: %v", res.StatusCode)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't read bytes of response: %v", err)
	}

	return bytes, nil
}

func GetLocationAreaByLocationBytes(locationArea string) ([]byte, error) {
	res, err := http.Get(BaseUrl + "location-area/" + locationArea)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't get location area by location for area: %s, %v", locationArea, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("couldn't get location area by location: %v", res.StatusCode)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't read bytes of response: %v", err)
	}

	return bytes, nil
}

func GetPokemonByName(name string) ([]byte, error) {
	res, err := http.Get(BaseUrl + "pokemon/" + name)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't get pokemon by name: %s, %v", name, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("couldn't get pokemon by name: %s, status code: %d", name, res.StatusCode)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't read bytes of response: %v", err)
	}

	return bytes, nil
}

package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseUrl string = "https://pokeapi.co/api/v2/"

type locationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

func getLocationAreaBytes(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't get location areas: %v", err)
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't read bytes of response: %v", err)
	}

	return bytes, nil
}

func unmarshalLocationArea(bytes []byte) (locationArea, error) {
	var locationArea locationArea
	if err := json.Unmarshal(bytes, &locationArea); err != nil {
		return locationArea, fmt.Errorf("couldn't unmarshal location areas: %v", err)
	}

	return locationArea, nil
}
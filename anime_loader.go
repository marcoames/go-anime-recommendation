package main

import (
	"encoding/json"
	"os"
)

// LoadAnimeData reads the JSON file and returns a slice of Anime structs.
func LoadAnimeData(filePath string) ([]Anime, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var animeList []Anime
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&animeList); err != nil {
		return nil, err
	}

	return animeList, nil
}

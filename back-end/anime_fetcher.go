package main

import (
	"fmt"
	"time"
)

func AnimeFetcher() {
	// Create a JikanClient instance.
	client := NewJikanClient()

	var allAnime []map[string]interface{}
	for page := 1; page <= 200; page++ { // Each page returns 50 results.
		fmt.Printf("\rFetching page %d...", page)
		topAnime, err := client.FetchTopAnime(page)
		if err != nil {
			fmt.Printf("Error fetching top anime: %v\n", err)
			return
		}
		allAnime = append(allAnime, topAnime...)

		// Add a delay between requests to avoid rate-limiting.
		time.Sleep(1 * time.Second)
	}

	// save data to a JSON file
	filename := "anime_data.json"
	if err := SaveToFile(allAnime, filename); err != nil {
		fmt.Printf("Error saving data to file: %v\n", err)
		return
	}

	fmt.Printf("Anime saved to %s\n", filename)
}

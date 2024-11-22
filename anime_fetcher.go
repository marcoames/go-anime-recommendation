package main

import (
	"fmt"
	"time"
)

func AnimeFetcher() {
	// Create a JikanClient instance.
	client := NewJikanClient()

	var allAnime []map[string]interface{}
	for page := 1; page <= 10; page++ { // Each page returns 50 results; 10 pages for top 500 anime.
		fmt.Printf("Fetching page %d...\n", page)
		topAnime, err := client.FetchTopAnime(page)
		if err != nil {
			fmt.Printf("Error fetching top anime: %v\n", err)
			return
		}
		allAnime = append(allAnime, topAnime...)

		// Add a delay between requests to avoid rate-limiting.
		time.Sleep(1 * time.Second)
	}

	// Save the top 500 anime to a file.
	filename := "top_500_anime.json"
	if err := SaveToFile(allAnime, filename); err != nil {
		fmt.Printf("Error saving data to file: %v\n", err)
		return
	}

	fmt.Printf("Top 500 anime saved to %s\n", filename)
}

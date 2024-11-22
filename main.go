package main

import (
	"fmt"
)

func main() {

	// Fetch top 500 anime from Jikan API
	// AnimeFetcher() // Capitalized function name

	// Load the saved anime data from file
	animeData, err := LoadAnimeData("top_500_anime.json") // Capitalized function name
	if err != nil {
		fmt.Printf("Error loading anime data: %v\n", err)
		return
	}

	// Print the loaded data
	// fmt.Println("Loaded Anime Data:")
	// for _, anime := range animeData {
	// 	fmt.Printf("Title: %-80s Score: %.2f\n", anime.Title, anime.Score)
	// }

	allFeatures := prepareFeatures(animeData)
	weights := map[string]float64{
		"score":       0.05,
		"popularity":  0.01,
		"genres":      1,
		"demographic": 0.1,
		"studios":     0.1,
	}

	// Encode features
	encodedFeatures := encodeFeatures(allFeatures, weights)

	// Print encoded features in a formatted way
	// fmt.Println("Encoded Features:")
	// for i, encodedFeature := range encodedFeatures {
	// 	fmt.Printf("Anime %d: %.2f\n", i+1, encodedFeature)
	// }

	// Anime index and k neighbors
	animeIndex := getAnimeIndex("Sousou no Frieren", animeData)
	k := 4

	// Get recommendations for the given animeIndex
	recommendations := findRecommendations(encodedFeatures, animeIndex, k)

	// Output the recommendations
	fmt.Println("Recommendations for Anime", animeData[animeIndex].Title)
	for _, recommendation := range recommendations {
		fmt.Printf("Recommended Anime: %v\n", animeData[recommendation].Title)
	}

}

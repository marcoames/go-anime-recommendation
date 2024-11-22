package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Load anime data from file
	animeData, err := LoadAnimeData("top_500_anime.json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading anime data: %v", err), http.StatusInternalServerError)
		return
	}

	// Get the anime title from query parameters
	animeTitle := r.URL.Query().Get("anime")
	if animeTitle == "" {
		http.Error(w, "Please provide an anime title using the 'anime' query parameter.", http.StatusBadRequest)
		return
	}

	// Prepare features
	allFeatures := prepareFeatures(animeData)
	weights := map[string]float64{
		"score":       0.05,
		"popularity":  0.01,
		"genres":      1,
		"demographic": 0.1,
		"studios":     0.1,
	}

	encodedFeatures := encodeFeatures(allFeatures, weights)

	// Find the index of the anime
	animeIndex := getAnimeIndex(animeTitle, animeData)
	if animeIndex == -1 {
		http.Error(w, fmt.Sprintf("Anime '%s' not found in the database.", animeTitle), http.StatusNotFound)
		return
	}

	// Get recommendations
	k := 4
	recommendations := findRecommendations(encodedFeatures, animeIndex, k)

	// Output the recommendations
	fmt.Fprintf(w, "Recommendations for Anime: %v\n", animeData[animeIndex].Title)
	for _, recommendation := range recommendations {
		fmt.Fprintf(w, "Recommended Anime: %v\n", animeData[recommendation].Title)
	}
}

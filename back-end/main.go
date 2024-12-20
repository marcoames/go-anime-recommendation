package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func main() {

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://animerecommendation-443313.rj.r.appspot.com", // Deployed frontend
			"http://localhost:3000",                               // Local React development
		},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"}, // Allow these HTTP methods
		AllowedHeaders: []string{"Content-Type"},           // Allow Content-Type header
	})

	// Fetch and save anime data
	// AnimeFetcher()

	// Handle the request with your handler function
	http.HandleFunc("/api/", handleRequest)

	// Wrap your handler with the CORS handler
	handler := corsHandler.Handler(http.DefaultServeMux)

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", handler)

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Load anime data from file
	animeData, err := LoadAnimeData("anime_data.json")
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
		"score":       0.3, // slightly higher for quality-based recommendations
		"popularity":  0.2, // moderate for trending recommendations
		"genres":      1,   // prioritize genres for better personalization
		"demographic": 0.1, // keeps demographic in consideration
		"studios":     0.2, // moderate studio influence
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
	// fmt.Fprintf(w, "Recommendations for Anime: %v\n", animeData[animeIndex].Title)
	// for _, recommendation := range recommendations {
	// 	fmt.Fprintf(w, "Recommended Anime: %v\n", animeData[recommendation].Title)
	// }

	response := struct {
		Anime           Anime   `json:"anime"`
		Recommendations []Anime `json:"recommendations"`
	}{
		Anime: animeData[animeIndex],
	}

	// Add the recommended anime to the response
	for _, recommendation := range recommendations {
		if recommendation >= 0 && recommendation < len(animeData) {
			response.Recommendations = append(response.Recommendations, animeData[recommendation])
		}
	}

	// Set the response header to 'application/json'
	w.Header().Set("Content-Type", "application/json")

	// Encode and return the response as JSON
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

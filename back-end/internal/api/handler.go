package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/marcoames/go-anime-recommendation/internal/anime"
	"github.com/marcoames/go-anime-recommendation/internal/recommendation"
)

type Handler struct {
	animeRepo *anime.Repository
}

func NewHandler(mongoURI string) (*Handler, error) {
	repo, err := anime.NewRepository(mongoURI)
	if err != nil {
		return nil, err
	}
	return &Handler{animeRepo: repo}, nil
}

func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// Load anime data from MongoDB
	animeData, err := h.animeRepo.LoadAnimeData()
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
	allFeatures := recommendation.PrepareFeatures(animeData)
	weights := map[string]float64{
		"score":       0.3, // slightly higher for quality-based recommendations
		"popularity":  0.2, // moderate for trending recommendations
		"genres":      1,   // prioritize genres for better personalization
		"demographic": 0.1, // keeps demographic in consideration
		"studios":     0.2, // moderate studio influence
	}

	encodedFeatures := recommendation.EncodeFeatures(allFeatures, weights)

	// Find the index of the anime
	animeIndex, err := h.animeRepo.GetAnimeIndex(animeTitle, animeData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get recommendations
	k := 4
	recommendations := recommendation.FindRecommendations(encodedFeatures, animeIndex, k)

	// Output the recommendations
	// fmt.Fprintf(w, "Recommendations for Anime: %v\n", animeData[animeIndex].Title)
	// for _, recommendation := range recommendations {
	// 	fmt.Fprintf(w, "Recommended Anime: %v\n", animeData[recommendation].Title)
	// }

	response := struct {
		Anime           anime.Anime   `json:"anime"`
		Recommendations []anime.Anime `json:"recommendations"`
	}{
		Anime: animeData[animeIndex],
	}

	// Add the recommended anime to the response
	for _, recommendation := range recommendations {
		if recommendation >= 0 && recommendation < len(animeData) {
			response.Recommendations = append(response.Recommendations, animeData[recommendation])
		}
	}

	w.Header().Set("Content-Type", "application/json")

	// Encode and return the response as JSON
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

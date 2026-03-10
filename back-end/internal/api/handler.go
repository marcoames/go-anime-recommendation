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
	weights := map[string]float64{
		"score":       0.1,
		"popularity":  0.15,
		"genres":      1.5,
		"demographic": 1.0,
		"themes":      1.2,
		"studios":     0.1,
	}

	encodedFeatures := recommendation.EncodeWithOneHot(animeData, weights)

	// Find the index of the anime
	animeIndex, err := h.animeRepo.GetAnimeIndex(animeTitle, animeData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get recommendations
	k := 4
	candidates := recommendation.FindRecommendations(encodedFeatures, animeIndex, 80)

	filteredRecommendations := anime.FilteredRecommendations(animeData, candidates, animeIndex, k)

	response := struct {
		Anime           anime.Anime   `json:"anime"`
		Recommendations []anime.Anime `json:"recommendations"`
	}{
		Anime: animeData[animeIndex],
	}

	for _, idx := range filteredRecommendations {
		if idx >= 0 && idx < len(animeData) {
			response.Recommendations = append(response.Recommendations, animeData[idx])
		}
	}

	w.Header().Set("Content-Type", "application/json")

	// Encode and return the response as JSON
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

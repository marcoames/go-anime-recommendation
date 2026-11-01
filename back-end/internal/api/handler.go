package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/marcoames/go-anime-recommendation/internal/anime"
	"github.com/marcoames/go-anime-recommendation/internal/recommendation"
)

type Handler struct {
	animeRepo animeRepository
	mu        sync.RWMutex
	animeData []anime.Anime
	features  [][]float64
	loaded    bool
}

type animeRepository interface {
	LoadAnimeData() ([]anime.Anime, error)
	GetAnimeIndex(string, []anime.Anime) (int, error)
	SearchAnimeTitles(string, int64) ([]string, error)
}

func NewHandler(repo animeRepository) *Handler {
	return &Handler{animeRepo: repo}
}

func (h *Handler) loadAnimeData() ([]anime.Anime, [][]float64, error) {
	h.mu.RLock()
	if h.loaded {
		data, features := h.animeData, h.features
		h.mu.RUnlock()
		return data, features, nil
	}
	h.mu.RUnlock()

	h.mu.Lock()
	defer h.mu.Unlock()
	if h.loaded {
		return h.animeData, h.features, nil
	}

	data, err := h.animeRepo.LoadAnimeData()
	if err != nil {
		return nil, nil, err
	}
	weights := map[string]float64{
		"score":       0.3,
		"popularity":  0.2,
		"genres":      1,
		"demographic": 0.15,
		"studios":     0.1,
	}
	h.animeData = data
	h.features = recommendation.EncodeFeatures(recommendation.PrepareFeatures(data), weights)
	h.loaded = true
	return h.animeData, h.features, nil
}

func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	animeData, encodedFeatures, err := h.loadAnimeData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading anime data: %v", err), http.StatusInternalServerError)
		return
	}

	animeTitle := r.URL.Query().Get("anime")
	if animeTitle == "" {
		http.Error(w, "Please provide an anime title using the 'anime' query parameter.", http.StatusBadRequest)
		return
	}

	animeIndex, err := h.animeRepo.GetAnimeIndex(animeTitle, animeData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	k := 6
	recommendations := recommendation.FindDiverseRecommendations(encodedFeatures, animeData, animeIndex, k)

	response := struct {
		Anime           anime.Anime   `json:"anime"`
		Recommendations []anime.Anime `json:"recommendations"`
	}{
		Anime: animeData[animeIndex],
	}

	for _, recommendation := range recommendations {
		if recommendation >= 0 && recommendation < len(animeData) {
			response.Recommendations = append(response.Recommendations, animeData[recommendation])
		}
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleSuggestions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		err := json.NewEncoder(w).Encode(map[string][]string{
			"suggestions": {},
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		}
		return
	}

	suggestions, err := h.animeRepo.SearchAnimeTitles(query, 8)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading suggestions: %v", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string][]string{
		"suggestions": suggestions,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

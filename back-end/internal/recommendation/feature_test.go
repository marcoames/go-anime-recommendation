package recommendation

import (
	"testing"

	"github.com/marcoames/go-anime-recommendation/internal/anime"
)

func TestBuildVocabulary(t *testing.T) {
	animeData := []anime.Anime{
		{
			Genres:       []anime.Genre{{Name: "Action"}, {Name: "Comedy"}},
			Demographics: []anime.Info{{Name: "Shounen"}},
			Studios:      []anime.Info{{Name: "Studio A"}},
			Themes:       []anime.Info{{Name: "Military"}},
		},
		{
			Genres:       []anime.Genre{{Name: "Action"}, {Name: "Drama"}},
			Demographics: []anime.Info{{Name: "Seinen"}},
			Studios:      []anime.Info{{Name: "Studio B"}},
			Themes:       []anime.Info{{Name: "Military"}, {Name: "Isekai"}},
		},
	}

	genreIdx, demoIdx, studioIdx, themeIdx := BuildVocabulary(animeData)

	if len(genreIdx) != 3 { // Action, Comedy, Drama
		t.Errorf("expected 3 genres, got %d", len(genreIdx))
	}
	if len(demoIdx) != 2 { // Shounen, Seinen
		t.Errorf("expected 2 demographics, got %d", len(demoIdx))
	}
	if len(studioIdx) != 2 { // Studio A, Studio B
		t.Errorf("expected 2 studios, got %d", len(studioIdx))
	}
	if len(themeIdx) != 2 { // Military, Isekai
		t.Errorf("expected 2 themes, got %d", len(themeIdx))
	}
}

func TestEncodeWithOneHot(t *testing.T) {
	animeData := []anime.Anime{
		{
			Score:        8.5,
			Popularity:   100,
			Genres:       []anime.Genre{{Name: "Action"}},
			Demographics: []anime.Info{},
			Studios:      []anime.Info{{Name: "Studio A"}},
		},
		{
			Score:        7.0,
			Popularity:   5000,
			Genres:       []anime.Genre{{Name: "Action"}, {Name: "Comedy"}},
			Demographics: []anime.Info{},
			Studios:      []anime.Info{{Name: "Studio A"}},
		},
	}

	weights := map[string]float64{
		"score":       1.0,
		"popularity":  1.0,
		"genres":      1.0,
		"demographic": 1.0,
		"studios":     1.0,
		"themes":      1.0,
	}

	encoded := EncodeWithOneHot(animeData, weights)

	if len(encoded) != 2 {
		t.Fatalf("expected 2 encoded vectors, got %d", len(encoded))
	}

	// Both vectors should have the same length
	if len(encoded[0]) != len(encoded[1]) {
		t.Error("encoded vectors have different lengths")
	}

	// First anime should have higher score component
	if encoded[0][0] <= encoded[1][0] {
		t.Error("expected first anime to have higher score")
	}
}

func TestSimilarGenresRankHigher(t *testing.T) {
	// Anime 0: Action, Comedy
	// Anime 1: Action, Comedy (should be most similar to 0)
	// Anime 2: Horror, Romance (should be least similar to 0)
	animeData := []anime.Anime{
		{Score: 8.0, Popularity: 100, Genres: []anime.Genre{{Name: "Action"}, {Name: "Comedy"}}},
		{Score: 8.0, Popularity: 100, Genres: []anime.Genre{{Name: "Action"}, {Name: "Comedy"}}},
		{Score: 8.0, Popularity: 100, Genres: []anime.Genre{{Name: "Horror"}, {Name: "Romance"}}},
	}

	weights := map[string]float64{
		"score":       0.1,
		"popularity":  0.1,
		"genres":      1.0,
		"demographic": 0.1,
		"studios":     0.1,
		"themes":      1.0,
	}

	encoded := EncodeWithOneHot(animeData, weights)
	recommendations := FindRecommendations(encoded, 0, 2)

	// Anime 1 should be the top recommendation (same genres)
	if recommendations[0] != 1 {
		t.Errorf("expected anime 1 as top recommendation, got %d", recommendations[0])
	}
}

func TestWeightConfigurations(t *testing.T) {
	// Create test data that mimics real anime
	animeData := []anime.Anime{
		// Index 0: Naruto-like (Shounen, Action, Adventure)
		{Title: "Query Anime", Score: 8.0, Popularity: 100,
			Genres:       []anime.Genre{{Name: "Action"}, {Name: "Adventure"}},
			Demographics: []anime.Info{{Name: "Shounen"}},
			Studios:      []anime.Info{{Name: "Studio A"}}},
		// Index 1: Similar shounen (should rank high)
		{Title: "Similar Shounen", Score: 8.5, Popularity: 150,
			Genres:       []anime.Genre{{Name: "Action"}, {Name: "Adventure"}},
			Demographics: []anime.Info{{Name: "Shounen"}},
			Studios:      []anime.Info{{Name: "Studio B"}}},
		// Index 2: Same studio, different genre
		{Title: "Same Studio", Score: 7.0, Popularity: 500,
			Genres:       []anime.Genre{{Name: "Romance"}},
			Demographics: []anime.Info{{Name: "Shoujo"}},
			Studios:      []anime.Info{{Name: "Studio A"}}},
		// Index 3: Different everything
		{Title: "Completely Different", Score: 9.0, Popularity: 50,
			Genres:       []anime.Genre{{Name: "Horror"}, {Name: "Mystery"}},
			Demographics: []anime.Info{{Name: "Seinen"}},
			Studios:      []anime.Info{{Name: "Studio C"}}},
	}

	testCases := []struct {
		name          string
		weights       map[string]float64
		expectedFirst int // Expected top recommendation index
	}{
		{
			name: "Genre Heavy - should recommend similar genres",
			weights: map[string]float64{
				"score": 0.1, "popularity": 0.1, "genres": 2.0, "demographic": 0.1, "studios": 0.1, "themes": 1.0,
			},
			expectedFirst: 1, // Similar Shounen
		},
		{
			name: "Studio Heavy - should recommend same studio",
			weights: map[string]float64{
				"score": 0.1, "popularity": 0.1, "genres": 0.1, "demographic": 0.1, "studios": 2.0, "themes": 1.0,
			},
			expectedFirst: 2, // Same Studio
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded := EncodeWithOneHot(animeData, tc.weights)
			recommendations := FindRecommendations(encoded, 0, 3)

			if recommendations[0] != tc.expectedFirst {
				t.Errorf("expected first recommendation to be %d (%s), got %d (%s)",
					tc.expectedFirst, animeData[tc.expectedFirst].Title,
					recommendations[0], animeData[recommendations[0]].Title)
			}
		})
	}
}

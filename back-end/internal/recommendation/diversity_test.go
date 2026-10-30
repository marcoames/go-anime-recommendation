package recommendation

import (
	"testing"

	"github.com/marcoames/go-anime-recommendation/internal/anime"
)

func TestFindDiverseRecommendationsAvoidsQueryFranchise(t *testing.T) {
	data := []anime.Anime{
		{Title: "Demon Slayer: Kimetsu no Yaiba", TitleSynonyms: []string{"Kimetsu no Yaiba"}, Genres: []anime.Genre{{Name: "Action"}}},
		{Title: "Kimetsu no Yaiba: Entertainment District Arc", TitleSynonyms: []string{"Demon Slayer"}, Genres: []anime.Genre{{Name: "Action"}}},
		{Title: "Jujutsu Kaisen", Genres: []anime.Genre{{Name: "Action"}}},
		{Title: "Attack on Titan", Genres: []anime.Genre{{Name: "Action"}}},
	}
	features := [][]float64{{1, 0}, {1, 0}, {0.95, 0}, {0.9, 0}}

	got := FindDiverseRecommendations(features, data, 0, 2)
	if len(got) != 2 || got[0] == 1 || got[1] == 1 {
		t.Fatalf("expected unrelated recommendations, got %v", got)
	}
}

func TestFindDiverseRecommendationsLimitsDuplicateFranchises(t *testing.T) {
	data := []anime.Anime{
		{Title: "Query", Genres: []anime.Genre{{Name: "Action"}}},
		{Title: "One Piece", Genres: []anime.Genre{{Name: "Action"}}},
		{Title: "One Piece Film: Red", Genres: []anime.Genre{{Name: "Action"}}},
		{Title: "Naruto", Genres: []anime.Genre{{Name: "Action"}}},
	}
	features := [][]float64{{1, 0}, {0.99, 0}, {0.98, 0}, {0.8, 0}}
	if !sharesFranchise(franchiseKeys(data[1]), franchiseKeys(data[2])) {
		t.Fatalf("expected franchise keys to match: %#v != %#v", franchiseKeys(data[1]), franchiseKeys(data[2]))
	}

	got := FindDiverseRecommendations(features, data, 0, 3)
	if len(got) != 2 {
		t.Fatalf("expected one result per franchise, got %v", got)
	}
}

func TestFindDiverseRecommendationsUsesFranchiseID(t *testing.T) {
	data := []anime.Anime{
		{Title: "Query", FranchiseID: 10},
		{Title: "Different Title", FranchiseID: 10},
		{Title: "Another Anime", FranchiseID: 20},
		{Title: "Third Anime", FranchiseID: 30},
	}
	features := [][]float64{{1, 0}, {0.99, 0}, {0.8, 0}, {0.7, 0}}

	got := FindDiverseRecommendations(features, data, 0, 2)
	if len(got) != 2 || got[0] != 2 || got[1] != 3 {
		t.Fatalf("expected different franchises, got %v", got)
	}
}

func TestFindDiverseRecommendationsRejectsInvalidInput(t *testing.T) {
	features := [][]float64{{1, 0}}
	data := []anime.Anime{{Title: "Query"}}

	if got := FindDiverseRecommendations(features, data, -1, 1); got != nil {
		t.Fatalf("expected nil for invalid query index, got %v", got)
	}
	if got := FindDiverseRecommendations(features, data, 0, 0); got != nil {
		t.Fatalf("expected nil for non-positive k, got %v", got)
	}
}

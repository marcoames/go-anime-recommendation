package recommendation

import (
	"testing"

	"github.com/marcoames/go-anime-recommendation/internal/anime"
)

func TestEncodeFeaturesUsesCategoryIdentity(t *testing.T) {
	data := []anime.Anime{
		{Score: 8, Popularity: 10, Genres: []anime.Genre{{Name: "Action"}}},
		{Score: 8, Popularity: 20, Genres: []anime.Genre{{Name: "Romance"}}},
		{Score: 8, Popularity: 30, Genres: []anime.Genre{{Name: "Action"}}},
	}

	features := EncodeFeatures(PrepareFeatures(data), map[string]float64{
		"score": 1, "popularity": 1, "genres": 1,
	})

	if len(features) != 3 || len(features[0]) != len(features[1]) {
		t.Fatalf("unexpected encoded feature dimensions: %#v", features)
	}
	if features[0][2] != 1 || features[2][2] != 1 || features[1][2] != 0 {
		t.Fatalf("genre identity was not preserved: %#v", features)
	}
}

func TestEncodeFeaturesNormalizesPopularity(t *testing.T) {
	features := EncodeFeatures([]Feature{
		{Score: 10, Popularity: 10},
		{Score: 10, Popularity: 100},
	}, map[string]float64{"score": 1, "popularity": 1})

	if features[0][1] <= features[1][1] {
		t.Fatalf("better popularity rank should have higher normalized value: %#v", features)
	}
}

func TestEncodeFeaturesDoesNotRewardMissingPopularity(t *testing.T) {
	features := EncodeFeatures([]Feature{
		{Score: 10, Popularity: 0},
		{Score: 10, Popularity: 100},
	}, map[string]float64{"score": 1, "popularity": 1})

	if features[0][1] != 0 {
		t.Fatalf("missing popularity should not receive a positive score: %#v", features)
	}
}

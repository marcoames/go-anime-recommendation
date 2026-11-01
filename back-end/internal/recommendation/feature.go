package recommendation

import (
	"sort"

	"github.com/marcoames/go-anime-recommendation/internal/anime"
)

type Feature struct {
	Score        float64
	Popularity   int
	Genres       []string
	Demographics []string
	Studios      []string
}

func PrepareFeatures(animeData []anime.Anime) []Feature {
	features := make([]Feature, 0, len(animeData))
	for _, item := range animeData {
		feature := Feature{
			Score:      item.Score,
			Popularity: item.Popularity,
		}
		for _, genre := range item.Genres {
			feature.Genres = append(feature.Genres, genre.Name)
		}
		for _, demographic := range item.Demographics {
			feature.Demographics = append(feature.Demographics, demographic.Name)
		}
		for _, studio := range item.Studios {
			feature.Studios = append(feature.Studios, studio.Name)
		}
		features = append(features, feature)
	}
	return features
}

func EncodeFeatures(features []Feature, weights map[string]float64) [][]float64 {
	maxPopularity := 1
	for _, feature := range features {
		if feature.Popularity > maxPopularity {
			maxPopularity = feature.Popularity
		}
	}

	genres := collectCategories(features, func(feature Feature) []string { return feature.Genres })
	demographics := collectCategories(features, func(feature Feature) []string { return feature.Demographics })
	studios := collectCategories(features, func(feature Feature) []string { return feature.Studios })

	encoded := make([][]float64, 0, len(features))
	for _, feature := range features {
		vector := []float64{
			(feature.Score / 10) * weights["score"],
			(1 - float64(feature.Popularity)/float64(maxPopularity)) * weights["popularity"],
		}
		vector = append(vector, encodeCategories(feature.Genres, genres, weights["genres"])...)
		vector = append(vector, encodeCategories(feature.Demographics, demographics, weights["demographic"])...)
		vector = append(vector, encodeCategories(feature.Studios, studios, weights["studios"])...)
		encoded = append(encoded, vector)
	}
	return encoded
}

func collectCategories(features []Feature, selector func(Feature) []string) []string {
	seen := make(map[string]struct{})
	for _, feature := range features {
		for _, category := range selector(feature) {
			if category != "" {
				seen[category] = struct{}{}
			}
		}
	}
	categories := make([]string, 0, len(seen))
	for category := range seen {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

func encodeCategories(values, categories []string, weight float64) []float64 {
	valueSet := make(map[string]struct{}, len(values))
	for _, value := range values {
		valueSet[value] = struct{}{}
	}
	encoded := make([]float64, len(categories))
	for i, category := range categories {
		if _, ok := valueSet[category]; ok {
			encoded[i] = weight
		}
	}
	return encoded
}

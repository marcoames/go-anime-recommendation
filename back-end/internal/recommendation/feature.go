package recommendation

import (
	"math"

	"github.com/marcoames/go-anime-recommendation/internal/anime"
)

// BuildVocabulary extracts all unique genres, demographics, studios, and themes
func BuildVocabulary(animeData []anime.Anime) (genreIndex, demoIndex, studioIndex, themeIndex map[string]int) {
	genreIndex = make(map[string]int)
	demoIndex = make(map[string]int)
	studioIndex = make(map[string]int)
	themeIndex = make(map[string]int)

	for _, a := range animeData {
		for _, g := range a.Genres {
			if _, exists := genreIndex[g.Name]; !exists {
				genreIndex[g.Name] = len(genreIndex)
			}
		}
		for _, d := range a.Demographics {
			if _, exists := demoIndex[d.Name]; !exists {
				demoIndex[d.Name] = len(demoIndex)
			}
		}
		for _, s := range a.Studios {
			if _, exists := studioIndex[s.Name]; !exists {
				studioIndex[s.Name] = len(studioIndex)
			}
		}
		for _, t := range a.Themes {
			if _, exists := themeIndex[t.Name]; !exists {
				themeIndex[t.Name] = len(themeIndex)
			}
		}
	}
	return
}

// EncodeWithOneHot creates feature vectors with one-hot encoded categorical features
func EncodeWithOneHot(animeData []anime.Anime, weights map[string]float64) [][]float64 {
	genreIdx, demoIdx, studioIdx, themeIdx := BuildVocabulary(animeData)
	numGenres := len(genreIdx)
	numDemos := len(demoIdx)
	numStudios := len(studioIdx)
	numThemes := len(themeIdx)

	// Find max popularity for normalization
	maxPop := 1.0
	for _, a := range animeData {
		if float64(a.Popularity) > maxPop {
			maxPop = float64(a.Popularity)
		}
	}

	// Feature vector: score + popularity + genres + demographics + studios + themes
	vectorLen := 2 + numGenres + numDemos + numStudios + numThemes
	encoded := make([][]float64, len(animeData))

	for i, a := range animeData {
		vec := make([]float64, vectorLen)

		// Normalized score (0-1)
		vec[0] = (a.Score / 10.0) * weights["score"]

		// Normalized popularity using log scale for better distribution
		normalizedPop := 1.0 - (math.Log1p(float64(a.Popularity)) / math.Log1p(maxPop))
		if normalizedPop < 0 {
			normalizedPop = 0
		}
		vec[1] = normalizedPop * weights["popularity"]

		// One-hot encode genres
		offset := 2
		for _, g := range a.Genres {
			idx := genreIdx[g.Name]
			vec[offset+idx] = weights["genres"]
		}

		// One-hot encode demographics
		offset += numGenres
		for _, d := range a.Demographics {
			idx := demoIdx[d.Name]
			vec[offset+idx] = weights["demographic"]
		}

		// One-hot encode studios
		offset += numDemos
		for _, s := range a.Studios {
			idx := studioIdx[s.Name]
			vec[offset+idx] = weights["studios"]
		}

		// One-hot encode themes
		offset += numStudios
		for _, t := range a.Themes {
			idx := themeIdx[t.Name]
			vec[offset+idx] = weights["themes"]
		}

		encoded[i] = vec
	}
	return encoded
}

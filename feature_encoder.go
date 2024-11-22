package main

import (
	"strings"
)

func prepareFeatures(animeData []Anime) [][]interface{} {
	var allFeatures [][]interface{}

	// Extract the individual features
	for i := 0; i < len(animeData); i++ {

		score := animeData[i].Score
		popularity := animeData[i].Popularity

		// Extract genre names
		genreNames := []string{}
		for _, genre := range animeData[i].Genres {
			genreNames = append(genreNames, genre.Name)
		}

		// Extract demographic names
		demographicNames := []string{}
		for _, demographic := range animeData[i].Demographics {
			demographicNames = append(demographicNames, demographic.Name)
		}

		// Extract studio names
		studioNames := []string{}
		for _, studio := range animeData[i].Studios {
			studioNames = append(studioNames, studio.Name)
		}

		// Combine all the features into a single slice
		features := []interface{}{
			score,
			popularity,
			strings.Join(genreNames, ", "),
			strings.Join(demographicNames, ", "),
			strings.Join(studioNames, ", "),
		}

		// Append the features to the list of all features
		allFeatures = append(allFeatures, features)
	}

	return allFeatures
}

func encodeFeatures(allFeatures [][]interface{}, weights map[string]float64) [][]float64 {
	encodedFeatures := [][]float64{} // Initialize a slice of slices

	for _, features := range allFeatures {
		// Extract features from each anime
		score := features[0].(float64)
		popularity := features[1].(int)
		genres := features[2].(string)
		demographics := features[3].(string)
		studios := features[4].(string)

		// Encode each feature by multiplying with its weight
		scoreWeighted := score * weights["score"]
		popularityWeighted := float64(popularity) * weights["popularity"]
		genresWeighted := float64(len(strings.Split(genres, ", "))) * weights["genres"]
		demographicsWeighted := float64(len(strings.Split(demographics, ", "))) * weights["demographic"]
		studiosWeighted := float64(len(strings.Split(studios, ", "))) * weights["studios"]

		// Combine the weighted features into a slice for the current anime
		encodedAnime := []float64{
			scoreWeighted,
			popularityWeighted,
			genresWeighted,
			demographicsWeighted,
			studiosWeighted,
		}

		// Append the encoded anime features to the final encodedFeatures slice
		encodedFeatures = append(encodedFeatures, encodedAnime)
	}

	return encodedFeatures
}

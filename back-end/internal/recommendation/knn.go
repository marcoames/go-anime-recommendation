package recommendation

import (
	"math"
	"sort"
)

func cosineSimilarity(v1, v2 []float64) float64 {
	var dotProduct, magV1, magV2 float64
	for i := range v1 {
		dotProduct += v1[i] * v2[i]
		magV1 += v1[i] * v1[i]
		magV2 += v2[i] * v2[i]
	}
	magV1 = math.Sqrt(magV1)
	magV2 = math.Sqrt(magV2)
	if magV1 == 0 || magV2 == 0 {
		return 0 // Avoid division by zero
	}
	return dotProduct / (magV1 * magV2)
}

func FindRecommendations(features [][]float64, queryIndex int, k int) []int {
	if queryIndex < 0 || queryIndex >= len(features) || k <= 0 {
		return nil
	}
	query := features[queryIndex]
	type neighbor struct {
		index    int
		distance float64
	}

	var neighbors []neighbor

	for i, feature := range features {
		if i == queryIndex {
			continue
		}
		if len(feature) != len(query) {
			continue
		}
		similarity := cosineSimilarity(query, feature)
		neighbors = append(neighbors, neighbor{index: i, distance: similarity})
	}

	sort.Slice(neighbors, func(i, j int) bool {
		return neighbors[i].distance > neighbors[j].distance
	})

	topK := make([]int, 0, k)
	for i := 0; i < k && i < len(neighbors); i++ {
		topK = append(topK, neighbors[i].index)
	}

	return topK
}

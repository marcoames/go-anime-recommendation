package recommendation

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/marcoames/go-anime-recommendation/internal/anime"
)

const (
	relevanceWeight = 0.9
	relevanceFloor  = 0.75
	candidatePool   = 50
)

var franchiseSuffixPattern = regexp.MustCompile(`\b(19|20)\d{2}\b|\b(\d+)(st|nd|rd|th)?\s+season\b|\b(final|part|cour|movie|film|ova|ona|special|recap|arc)\b.*$`)

type recommendationCandidate struct {
	index     int
	relevance float64
}

// FindDiverseRecommendations balances relevance and diversity while avoiding
// repeated franchises when sufficiently relevant alternatives exist.
func FindDiverseRecommendations(features [][]float64, animeData []anime.Anime, queryIndex, k int) []int {
	if queryIndex < 0 || queryIndex >= len(features) || queryIndex >= len(animeData) || k <= 0 {
		return nil
	}

	query := features[queryIndex]
	queryFranchises := franchiseKeys(animeData[queryIndex])
	candidates := make([]recommendationCandidate, 0, len(features)-1)
	for index, feature := range features {
		if index == queryIndex || len(feature) != len(query) {
			continue
		}
		candidates = append(candidates, recommendationCandidate{
			index:     index,
			relevance: cosineSimilarity(query, feature),
		})
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].relevance > candidates[j].relevance
	})
	if len(candidates) == 0 {
		return nil
	}

	poolSize := candidatePool
	if poolSize > len(candidates) {
		poolSize = len(candidates)
	}
	candidates = candidates[:poolSize]
	minRelevance := candidates[0].relevance * relevanceFloor
	relevantCandidates := make([]recommendationCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.relevance >= minRelevance {
			relevantCandidates = append(relevantCandidates, candidate)
		}
	}

	selected := make([]int, 0, k)
	selectedSet := make(map[int]struct{}, k)
	franchises := make(map[string]struct{}, k)
	for len(selected) < k && len(selected) < len(relevantCandidates) {
		bestIndex := -1
		bestScore := -1.0
		for _, candidate := range relevantCandidates {
			if _, ok := selectedSet[candidate.index]; ok {
				continue
			}
			candidateFranchises := franchiseKeys(animeData[candidate.index])
			if sharesFranchise(queryFranchises, candidateFranchises) {
				continue
			}
			if hasSelectedFranchise(candidateFranchises, franchises) {
				continue
			}

			score := mmrScore(candidate, selected, features)
			if score > bestScore {
				bestIndex = candidate.index
				bestScore = score
			}
		}
		if bestIndex == -1 {
			break
		}
		selected = append(selected, bestIndex)
		selectedSet[bestIndex] = struct{}{}
		addFranchises(franchises, franchiseKeys(animeData[bestIndex]))
	}

	for len(selected) < k {
		bestIndex := -1
		bestScore := -1.0
		for _, candidate := range relevantCandidates {
			if _, ok := selectedSet[candidate.index]; ok {
				continue
			}
			if hasSelectedFranchise(franchiseKeys(animeData[candidate.index]), franchises) {
				continue
			}
			score := mmrScore(candidate, selected, features)
			if score > bestScore {
				bestIndex = candidate.index
				bestScore = score
			}
		}
		if bestIndex == -1 {
			break
		}
		selected = append(selected, bestIndex)
		selectedSet[bestIndex] = struct{}{}
		addFranchises(franchises, franchiseKeys(animeData[bestIndex]))
	}
	return selected
}

func mmrScore(candidate recommendationCandidate, selected []int, features [][]float64) float64 {
	if len(selected) == 0 {
		return candidate.relevance
	}
	maxRedundancy := 0.0
	for _, selectedIndex := range selected {
		redundancy := cosineSimilarity(features[candidate.index], features[selectedIndex])
		if redundancy > maxRedundancy {
			maxRedundancy = redundancy
		}
	}
	return relevanceWeight*candidate.relevance - (1-relevanceWeight)*maxRedundancy
}

func franchiseKeys(item anime.Anime) []string {
	keys := []string{
		normalizeFranchiseTitle(item.Title),
		normalizeFranchiseTitle(item.TitleEnglish),
		normalizeFranchiseTitle(item.TitleJapanese),
	}
	for _, synonym := range item.TitleSynonyms {
		keys = append(keys, normalizeFranchiseTitle(synonym))
	}
	if item.FranchiseID > 0 {
		keys = append(keys, "mal-franchise-"+strconv.Itoa(item.FranchiseID))
	}
	unique := make([]string, 0, len(keys))
	seen := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		if key != "" {
			if _, ok := seen[key]; !ok {
				seen[key] = struct{}{}
				unique = append(unique, key)
			}
		}
	}
	return unique
}

func sharesFranchise(left, right []string) bool {
	for _, leftKey := range left {
		for _, rightKey := range right {
			if leftKey == rightKey {
				return true
			}
		}
	}
	return false
}

func hasSelectedFranchise(keys []string, selected map[string]struct{}) bool {
	for _, key := range keys {
		if _, ok := selected[key]; ok {
			return true
		}
	}
	return false
}

func addFranchises(selected map[string]struct{}, keys []string) {
	for _, key := range keys {
		selected[key] = struct{}{}
	}
}

func normalizeFranchiseTitle(title string) string {
	title = strings.ToLower(strings.TrimSpace(title))
	if title == "" {
		return ""
	}
	title = strings.ReplaceAll(title, "-", " ")
	title = franchiseSuffixPattern.ReplaceAllString(title, "")
	if colon := strings.Index(title, ":"); colon >= 0 {
		title = title[:colon]
	}
	title = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(title, " ")
	return strings.Join(strings.Fields(title), " ")
}

func franchiseKey(item anime.Anime) string {
	keys := franchiseKeys(item)
	if len(keys) == 0 {
		return ""
	}
	return keys[0]
}

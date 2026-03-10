package recommendation

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/marcoames/go-anime-recommendation/internal/anime"
)

// expectedPair defines a query anime and a list of anime titles that
// should reasonably appear in the top-k recommendations.
type expectedPair struct {
	query    string
	expected []string // at least one of these should appear in top-k
}

// groundTruth returns hand-curated pairs based on common knowledge / MAL recommendations.
func groundTruth() []expectedPair {
	return []expectedPair{
		{
			query: "Naruto",
			expected: []string{
				"Bleach", "One Piece", "Hunter x Hunter", "Dragon Ball Z",
				"Fairy Tail", "Black Clover", "Naruto: Shippuuden",
				"Nanatsu no Taizai", "Magi",
			},
		},
		{
			query: "Death Note",
			expected: []string{
				"Code Geass", "Monster", "Psycho-Pass", "Steins;Gate",
				"Mirai Nikki", "Summertime Render", "Deadman Wonderland",
				"Tomodachi Game", "Classroom of the Elite",
				"Youkoso Jitsuryoku",
			},
		},
		{
			query: "Attack on Titan",
			expected: []string{
				"Vinland Saga", "Demon Slayer", "Fullmetal Alchemist",
				"Tokyo Ghoul", "Psycho-Pass", "86", "Kabaneri",
				"91 Days", "Madoka", "Phantom",
			},
		},
		{
			query: "Steins;Gate",
			expected: []string{
				"Re:Zero", "Erased", "Death Note", "Madoka Magica",
				"Psycho-Pass", "Steins;Gate 0", "Evangelion",
				"Guilty Crown", "Inuyashiki", "Classroom of the Elite",
				"Youkoso Jitsuryoku",
			},
		},
		{
			query: "Sword Art Online",
			expected: []string{
				"No Game No Life", "Log Horizon", "Overlord",
				"That Time I Got Reincarnated as a Slime",
				"Gate: Jieitai", "Dungeon ni Deai",
				"Solo Leveling", "Ore dake Level",
			},
		},
	}
}

// findIndex returns the index of an anime by title (case-insensitive, checks english title too).
// Prefers exact matches over partial matches.
func findIndex(animeData []anime.Anime, title string) int {
	// First pass: exact match
	for i, a := range animeData {
		if strings.EqualFold(a.Title, title) ||
			strings.EqualFold(a.TitleEnglish, title) {
			return i
		}
	}
	// Second pass: partial match, but only if the title lengths are similar
	lowerTitle := strings.ToLower(title)
	bestIdx := -1
	bestLen := int(^uint(0) >> 1) // max int
	for i, a := range animeData {
		lt := strings.ToLower(a.Title)
		lte := strings.ToLower(a.TitleEnglish)
		if strings.Contains(lt, lowerTitle) && len(a.Title) < bestLen {
			bestIdx = i
			bestLen = len(a.Title)
		} else if strings.Contains(lte, lowerTitle) && len(a.TitleEnglish) < bestLen {
			bestIdx = i
			bestLen = len(a.TitleEnglish)
		}
	}
	return bestIdx
}

// containsAny checks if any of the expected titles appear in the recommendation indices
func containsAny(animeData []anime.Anime, recIndices []int, expectedTitles []string) (bool, string) {
	for _, idx := range recIndices {
		for _, exp := range expectedTitles {
			if strings.Contains(strings.ToLower(animeData[idx].Title), strings.ToLower(exp)) ||
				strings.Contains(strings.ToLower(animeData[idx].TitleEnglish), strings.ToLower(exp)) {
				return true, exp
			}
		}
	}
	return false, ""
}

// TestRecommendationQuality runs against data loaded from MongoDB.
func TestRecommendationQuality(t *testing.T) {
	_ = godotenv.Load("../../../.env")
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		t.Skip("Skipping: MONGODB_URI not set in environment or .env")
	}

	repo, err := anime.NewRepository(mongoURI)
	if err != nil {
		t.Skipf("Skipping validation test, cannot connect to MongoDB: %v", err)
	}

	animeData, err := repo.LoadAnimeData()
	if err != nil || len(animeData) == 0 {
		t.Skipf("Skipping validation test, no anime data available: %v", err)
	}

	weights := map[string]float64{
		"score":       0.1,
		"popularity":  0.15,
		"genres":      1.5,
		"demographic": 1.0,
		"themes":      1.2,
		"studios":     0.1,
	}

	encoded := EncodeWithOneHot(animeData, weights)
	pairs := groundTruth()
	k := 10

	hits := 0
	total := 0

	for _, pair := range pairs {
		queryIdx := findIndex(animeData, pair.query)
		if queryIdx == -1 {
			t.Logf("SKIP: '%s' not found in dataset", pair.query)
			continue
		}

		total++
		candidates := FindRecommendations(encoded, queryIdx, 80)
		filtered := anime.FilteredRecommendations(animeData, candidates, queryIdx, k)

		found, matchedTitle := containsAny(animeData, filtered, pair.expected)

		// Log all recommendations for inspection
		recTitles := make([]string, len(filtered))
		for i, idx := range filtered {
			recTitles[i] = animeData[idx].Title
		}

		if found {
			hits++
			t.Logf("HIT:  '%s' -> matched '%s' | recs: %v", pair.query, matchedTitle, recTitles)
		} else {
			t.Logf("MISS: '%s' -> expected one of %v | got: %v", pair.query, pair.expected, recTitles)
		}
	}

	if total == 0 {
		t.Skip("No test pairs could be evaluated")
	}

	hitRate := float64(hits) / float64(total) * 100
	t.Logf("\n=== VALIDATION RESULTS ===")
	t.Logf("Hit Rate: %.1f%% (%d/%d)", hitRate, hits, total)
	t.Logf("==========================")

	// Fail if hit rate is below threshold
	if hitRate < 60.0 {
		t.Errorf("Hit rate %.1f%% is below 60%% threshold. Recommendations need tuning.", hitRate)
	}
}

// TestPrintRecommendations is a helper to manually inspect recommendations.
// Run with: go test -run TestPrintRecommendations -v
func TestPrintRecommendations(t *testing.T) {
	_ = godotenv.Load("../../../.env")
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		t.Skip("Skipping: MONGODB_URI not set in environment or .env")
	}

	repo, err := anime.NewRepository(mongoURI)
	if err != nil {
		t.Skipf("Skipping: %v", err)
	}

	animeData, err := repo.LoadAnimeData()
	if err != nil || len(animeData) == 0 {
		t.Skipf("Skipping: no data")
	}

	weights := map[string]float64{
		"score":       0.1,
		"popularity":  0.15,
		"genres":      1.5,
		"demographic": 1.0,
		"themes":      1.2,
		"studios":     0.1,
	}

	encoded := EncodeWithOneHot(animeData, weights)

	queries := []string{"Naruto", "Death Note", "Attack on Titan", "Steins;Gate", "One Piece", "Sword Art Online"}

	k := 10

	for _, q := range queries {
		idx := findIndex(animeData, q)
		if idx == -1 {
			t.Logf("'%s' not found in dataset", q)
			continue
		}

		candidates := FindRecommendations(encoded, idx, 80)
		filtered := anime.FilteredRecommendations(animeData, candidates, idx, k)

		fmt.Printf("\n=== Recommendations for '%s' ===\n", animeData[idx].Title)
		for rank, recIdx := range filtered {
			a := animeData[recIdx]
			fmt.Printf("  %2d. %s (Score: %.2f, Genres: %v)\n", rank+1, a.Title, a.Score, a.Genres)
		}
	}
}

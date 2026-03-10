package anime

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(uri string) (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database("anime").Collection("animes")
	return &Repository{collection: collection}, nil
}

func (r *Repository) LoadAnimeData() ([]Anime, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var animes []Anime
	if err = cursor.All(ctx, &animes); err != nil {
		return nil, err
	}

	return animes, nil
}

func (r *Repository) SaveAnimeData(animes []Anime) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert animes to interface{} for bulk write
	var documents []interface{}
	for _, anime := range animes {
		documents = append(documents, anime)
	}

	// Clear existing data
	if err := r.collection.Drop(ctx); err != nil {
		return err
	}

	// Insert new data
	_, err := r.collection.InsertMany(ctx, documents)
	return err
}

func (r *Repository) GetAnimeIndex(animeTitle string, animeData []Anime) (int, error) {
	for i, anime := range animeData {
		if strings.EqualFold(anime.Title, animeTitle) ||
			strings.EqualFold(anime.TitleEnglish, animeTitle) ||
			strings.EqualFold(anime.TitleJapanese, animeTitle) {
			return i, nil
		}
		for _, synonym := range anime.TitleSynonyms {
			if strings.EqualFold(synonym, animeTitle) {
				return i, nil
			}
		}
	}
	return -1, fmt.Errorf("Anime '%s' not found", animeTitle)
}

// extractBaseTitle removes common sequel/season indicators from a title
func extractBaseTitle(title string) string {
	result := strings.ToLower(strings.TrimSpace(title))
	if result == "" {
		return ""
	}

	// Remove content in parentheses at the end, e.g., "Title (TV)"
	if idx := strings.LastIndex(result, "("); idx > 0 {
		result = strings.TrimSpace(result[:idx])
	}

	// Order matters: remove longer patterns first to avoid partial matches
	removals := []string{
		": the final season", ": final season",
		" the final season", " final season",
		" 2nd season", " 3rd season", " 4th season", " 5th season",
		" season 2", " season 3", " season 4", " season 5",
		" part 2", " part 3", " part 4", " part ii", " part iii", " part iv",
		" cour 2", " cour 3",
		" next generations",
		" shippuuden", " shippuden",
		" brotherhood",
		" zoku", " shin", " new",
		" movie", " ova", " ona", " special", " specials",
		" recap",
	}

	for _, r := range removals {
		result = strings.TrimSuffix(result, r)
	}

	// Remove trailing roman numerals
	for _, suffix := range []string{" ii", " iii", " iv", " v", " vi", " vii"} {
		if strings.HasSuffix(result, suffix) {
			result = result[:len(result)-len(suffix)]
			break
		}
	}

	// Remove trailing digits like " 2", " 3" only if preceded by a letter
	trimmed := strings.TrimSpace(result)
	if len(trimmed) > 2 {
		last := trimmed[len(trimmed)-1]
		secondLast := trimmed[len(trimmed)-2]
		if last >= '2' && last <= '9' && secondLast == ' ' {
			result = trimmed[:len(trimmed)-2]
		}
	}

	// Split on ":" and take first part only if the prefix is long enough
	if idx := strings.Index(result, ":"); idx > 5 {
		result = result[:idx]
	}

	// Split on " - " and take first part
	if idx := strings.Index(result, " - "); idx > 5 {
		result = result[:idx]
	}

	return strings.TrimSpace(result)
}

func FilteredRecommendations(animeData []Anime, candidates []int, sourceIndex int, k int) []int {
	source := animeData[sourceIndex]
	baseTitle := extractBaseTitle(source.Title)
	baseEnglish := extractBaseTitle(source.TitleEnglish)

	seen := make(map[string]bool)
	// Mark source base titles as seen
	if baseTitle != "" {
		seen[baseTitle] = true
	}
	if baseEnglish != "" {
		seen[baseEnglish] = true
	}

	var results []int
	for _, idx := range candidates {
		if idx < 0 || idx >= len(animeData) || idx == sourceIndex {
			continue
		}

		candidate := animeData[idx]
		candidateBase := extractBaseTitle(candidate.Title)
		candidateBaseEnglish := extractBaseTitle(candidate.TitleEnglish)

		// Skip if the candidate shares the same base title (sequel/prequel)
		if (candidateBase != "" && seen[candidateBase]) ||
			(candidateBaseEnglish != "" && seen[candidateBaseEnglish]) {
			continue
		}

		// Skip if we already recommended something with the same base title (avoid duplicates from same franchise)
		if candidateBase != "" {
			if seen[candidateBase] {
				continue
			}
		}

		// Prefer same type (TV->TV, Movie->Movie), but don't hard-block
		// This is handled by scoring, not filtering

		if candidateBase != "" {
			seen[candidateBase] = true
		}
		if candidateBaseEnglish != "" {
			seen[candidateBaseEnglish] = true
		}

		results = append(results, idx)
		if len(results) >= k {
			break
		}
	}
	return results
}

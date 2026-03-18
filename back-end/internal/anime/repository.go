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

func (r *Repository) SearchAnimeTitles(query string, limit int64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query = strings.TrimSpace(query)
	if query == "" {
		return []string{}, nil
	}

	filter := bson.M{
		"$or": []bson.M{
			{
				"title": bson.M{
					"$regex":   "^" + query,
					"$options": "i",
				},
			},
			{
				"title_english": bson.M{
					"$regex":   "^" + query,
					"$options": "i",
				},
			},
			{
				"title_japanese": bson.M{
					"$regex":   "^" + query,
					"$options": "i",
				},
			},
			{
				"title_synonyms": bson.M{
					"$regex":   "^" + query,
					"$options": "i",
				},
			},
		},
	}

	opts := options.Find().
		SetLimit(limit).
		SetProjection(bson.M{
			"title": 1,
			"_id":   0,
		})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type animeTitleResult struct {
		Title string `bson:"title"`
	}

	var results []animeTitleResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	titles := make([]string, 0, len(results))

	for _, result := range results {
		title := strings.TrimSpace(result.Title)
		if title == "" || seen[title] {
			continue
		}
		seen[title] = true
		titles = append(titles, title)
	}

	return titles, nil
}

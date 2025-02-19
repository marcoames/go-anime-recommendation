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

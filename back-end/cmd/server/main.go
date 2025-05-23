package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/marcoames/go-anime-recommendation/internal/anime"
	"github.com/marcoames/go-anime-recommendation/internal/api"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	fetch := flag.Bool("fetch", false, "Fetch anime data before starting server")
	flag.Parse()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("❌ MONGODB_URI environment variable is not set")
	}

	// Create MongoDB client with Stable API options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	// Ping to confirm connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("❌ Ping to MongoDB failed: %v", err)
	}

	log.Println("✅ Connected to MongoDB Atlas!")

	// Ensure disconnection on exit
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("❌ Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Initialize repository with MongoDB client
	repo, err := anime.NewRepositoryWithClient(client)
	if err != nil {
		log.Fatalf("❌ Failed to create repository: %v", err)
	}

	// Optional fetch step
	if *fetch {
		log.Println("📦 Fetching anime data...")
		if err := anime.FetchAndSaveAnime(repo); err != nil {
			log.Fatalf("❌ Failed to fetch anime: %v", err)
		}
		log.Println("✅ Fetch complete")
		return
	}

	// Create API handler with MongoDB client
	handler, err := api.NewHandlerWithClient(client)
	if err != nil {
		log.Fatalf("❌ Failed to create handler: %v", err)
	}

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://go-anime-recommendation-1.onrender.com",
			"http://localhost:3000",
		},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	http.HandleFunc("/api/", handler.HandleRequest)
	wrappedHandler := corsHandler.Handler(http.DefaultServeMux)

	log.Println("🚀 Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", wrappedHandler))
}

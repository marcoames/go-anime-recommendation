package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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

	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is required")
	}
	
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)
	
	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()
	
	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB!")

	// If fetch flag is set, fetch data and exit
	if *fetch {
		log.Println("Fetching anime data...")
		repo, err := anime.NewRepository(mongoURI)
		if err != nil {
			log.Fatalf("Failed to create repository: %v", err)
		}
		
		if err := anime.FetchAndSaveAnime(repo); err != nil {
			log.Fatalf("Failed to fetch anime: %v", err)
		}
		log.Println("Fetch complete")
		return
	}

	// Create API handler
	handler, err := api.NewHandler(mongoURI)
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}

	// Setup CORS with proper configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://go-anime-recommendation-1.onrender.com", // No trailing slash
			"http://localhost:3000",
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Setup routes
	http.HandleFunc("/api/", handler.HandleRequest)
	
	// Add a health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap with CORS
	wrappedHandler := corsHandler.Handler(http.DefaultServeMux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, wrappedHandler))
}

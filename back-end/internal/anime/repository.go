package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/marcoames/go-anime-recommendation/internal/anime"
	"github.com/marcoames/go-anime-recommendation/internal/api"
	"github.com/rs/cors"
)

func main() {
	fetch := flag.Bool("fetch", false, "Fetch anime data before starting server")
	flag.Parse()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is required")
	}

	repo, err := anime.NewRepository(mongoURI)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	if *fetch {
		log.Println("Fetching anime data...")
		if err := anime.FetchAndSaveAnime(repo); err != nil {
			log.Fatalf("Failed to fetch anime: %v", err)
		}
		log.Println("Fetch complete")
		return
	}

	handler, err := api.NewHandler(mongoURI)
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://go-anime-recommendation-1.onrender.com", // Fixed: removed trailing slash
			"http://localhost:3000", // Local React development
		},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"}, // Allow these HTTP methods
		AllowedHeaders: []string{"Content-Type"},           // Allow Content-Type header
	})

	http.HandleFunc("/api/", handler.HandleRequest)
	wrappedHandler := corsHandler.Handler(http.DefaultServeMux)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", wrappedHandler))
}

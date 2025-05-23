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

	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is required")
	}

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

	// Setup CORS 
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://go-anime-recommendation-1.onrender.com",
			"http://localhost:3000",
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Setup routes
	http.HandleFunc("/api/", handler.HandleRequest)
	
	// health check endpoint
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

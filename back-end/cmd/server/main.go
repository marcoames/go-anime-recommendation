package main

import (
	"context"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/marcoames/go-anime-recommendation/internal/anime"
	"github.com/marcoames/go-anime-recommendation/internal/api"
	"github.com/rs/cors"
)

func main() {
	// Load .env
	_ = godotenv.Load()

	logOutput := io.Writer(os.Stdout)
	if os.Getenv("LOG_COLOR") == "1" && os.Getenv("NO_COLOR") == "" {
		logOutput = colorWriter{output: os.Stdout}
	}
	logger := slog.New(slog.NewTextHandler(logOutput, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	fetch := flag.Bool("fetch", false, "Fetch anime data before starting server")
	flag.Parse()

	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		logger.Error("MONGODB_URI environment variable is required")
		os.Exit(1)
	}

	repo, err := anime.NewRepository(mongoURI)
	if err != nil {
		logger.Error("failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err = repo.Close(context.Background()); err != nil {
			logger.Error("failed to disconnect from MongoDB", "error", err)
		}
	}()
	logger.Info("connected to MongoDB")

	// If fetch flag is set, fetch data and exit
	if *fetch {
		logger.Info("starting anime data fetch")
		if err := anime.FetchAndSaveAnime(repo); err != nil {
			logger.Error("failed to fetch anime", "error", err)
			os.Exit(1)
		}
		logger.Info("anime data fetch completed")
		return
	}

	// Create API handler
	handler := api.NewHandler(repo)

	// Setup CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://go-anime-recommendation-1.onrender.com",
			"http://localhost:3000",
		},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{
			"Content-Type", "Authorization", "X-Requested-With",
		},
		ExposedHeaders:     []string{"Content-Length"},
		AllowCredentials:   true,
		OptionsPassthrough: false,
		Debug:              false,
	})

	// Setup routes
	http.HandleFunc("/api/", handler.HandleRequest)
	http.HandleFunc("/api/suggestions", handler.HandleSuggestions)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap with CORS
	wrappedHandler := requestLogger(logger, corsHandler.Handler(http.DefaultServeMux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           wrappedHandler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	logger.Info("server is running", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}

type colorWriter struct {
	output io.Writer
}

func (w colorWriter) Write(p []byte) (int, error) {
	message := string(p)
	levelStart := strings.Index(message, "level=")
	if levelStart == -1 {
		return w.writeWhite(message, len(p))
	}

	levelValueStart := levelStart + len("level=")
	levelValueEnd := levelValueStart + strings.IndexAny(message[levelValueStart:], " \t\r\n")
	if levelValueEnd < levelValueStart {
		levelValueEnd = len(message)
	}

	levelColor := map[string]string{
		"ERROR": "\x1b[31m",
		"WARN":  "\x1b[33m",
		"INFO":  "\x1b[32m",
		"DEBUG": "\x1b[90m",
	}[message[levelValueStart:levelValueEnd]]
	if levelColor == "" {
		return w.writeWhite(message, len(p))
	}

	colored := "\x1b[37m" +
		message[:levelValueStart] +
		levelColor +
		message[levelValueStart:levelValueEnd] +
		"\x1b[37m" +
		message[levelValueEnd:] +
		"\x1b[0m"
	if _, err := w.output.Write([]byte(colored)); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w colorWriter) writeWhite(message string, length int) (int, error) {
	if _, err := w.output.Write([]byte("\x1b[37m" + message + "\x1b[0m")); err != nil {
		return 0, err
	}
	return length, nil
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(body []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(body)
}

func requestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w}

		next.ServeHTTP(recorder, r)

		status := recorder.status
		if status == 0 {
			status = http.StatusOK
		}
		logger.Info("http request",
			"method", r.Method,
			"path", r.URL.RequestURI(),
			"status", status,
			"duration_ms", time.Since(start).Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
	})
}

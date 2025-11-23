package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"s3-music-streamer/internal/config"
	"s3-music-streamer/internal/database"
	"s3-music-streamer/internal/handlers"
	"s3-music-streamer/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	ctx := context.Background()
	s3Client, err := storage.NewS3Client(ctx, cfg.S3Region, cfg.S3Bucket, cfg.AWSAccessKey, cfg.AWSSecretKey)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	handler := handlers.New(db, s3Client)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/songs", handler.ListSongs)
		r.Post("/songs", handler.CreateSong)
		r.Post("/songs/upload", handler.UploadSong)
		r.Get("/songs/{id}", handler.GetSong)
		r.Delete("/songs/{id}", handler.DeleteSong)
		r.Get("/songs/{id}/stream", handler.StreamSong)
	})

	// Serve static files from Client/dist
	workDir, _ := os.Getwd()
	clientPath := filepath.Join(workDir, "..", "Client", "dist")

	// Check if dist directory exists
	if _, err := os.Stat(clientPath); err == nil {
		fileServer := http.FileServer(http.Dir(clientPath))
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			// Try to serve the file
			filePath := filepath.Join(clientPath, r.URL.Path)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				// File doesn't exist, serve index.html for client-side routing
				http.ServeFile(w, r, filepath.Join(clientPath, "index.html"))
				return
			}
			fileServer.ServeHTTP(w, r)
		})
		log.Printf("Serving static files from %s", clientPath)
	} else {
		log.Printf("Warning: Client dist directory not found at %s", clientPath)
	}

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

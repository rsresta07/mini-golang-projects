package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"url-shortener-go/db"
	"url-shortener-go/handlers"
)

func main() {
	// load env
	_ = godotenv.Load()

	if err := db.Connect(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	log.Println("Connected to PostgreSQL via GORM")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not specified
	}

	// register routes on the router
	router := chi.NewRouter()

	// temp test route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL shortener is up and running!"))
	})

	router.Get("/{code}", handlers.RedirectHandler())

	router.Post("/api/shorten", handlers.ShortenHandler())

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("Server starting on http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// shutdown - this makes you look pro lol
	// wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // block here until signal

	log.Println("Shutting down server...")

	// create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// this will close connections cleanly
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

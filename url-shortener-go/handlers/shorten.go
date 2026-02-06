package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"time"
	"url-shortener-go/db"
	"url-shortener-go/models"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const codeLength = 6

// Local random generator
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// generateCode generates a random code of length codeLength using the
// provided charset. It is used to generate short codes for URLs.
// The generated code is a string of length codeLength containing only
// characters from the charset.
func generateCode() string {
	b := make([]byte, codeLength)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// ShortenHandler is an HTTP handler that takes a URL as a request payload and returns a shortened URL.
// It generates a random short code of length codeLength using the provided charset, and checks if the generated code already exists in the database.
// If a collision occurs (i.e. the generated code already exists in the database), it retries up to 10 times to generate a new code.
// If all retries fail, it returns an error response with a 500 status.
// If the short URL is successfully created, it returns a JSON response with the shortened URL.
func ShortenHandler() http.HandlerFunc {
	baseURL := os.Getenv("BASE_URL")

	return func(w http.ResponseWriter, r *http.Request) {
		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		var shortCode string
		for i := 0; i < 10; i++ {
			shortCode = generateCode()
			var existing models.URL
			if err := db.DB.Where("short_code = ?", shortCode).First(&existing).Error; err != nil {
				break
			}
		}

		url := models.URL{ShortCode: shortCode, LongURL: req.URL}
		if err := db.DB.Create(&url).Error; err != nil {
			http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
			return
		}

		resp := ShortenResponse{
			ShortURL: baseURL + "/" + shortCode,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

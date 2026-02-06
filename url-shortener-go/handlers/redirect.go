package handlers

import (
	"net/http"
	"url-shortener-go/db"
	"url-shortener-go/models"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// RedirectHandler is a HTTP handler that redirects a request to the long URL associated with the given short code.
// If the short code does not exist in the database, it returns a 404 status.
// If an error occurs while querying the database or updating the click count, it returns a 500 status.
// If the request is successful, it redirects the user to the long URL associated with the given short code.
func RedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		var url models.URL

		if err := db.DB.Where("short_code = ?", code).First(&url).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.NotFound(w, r)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		db.DB.Model(&url).Update("clicks", url.Clicks+1)
		http.Redirect(w, r, url.LongURL, http.StatusFound)
	}
}

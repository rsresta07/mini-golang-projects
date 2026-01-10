package handlers

import (
	"net/http"
	"url-shortener-go/db"
	"url-shortener-go/models"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

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

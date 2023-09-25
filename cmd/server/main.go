package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"web-scraper/pkg/scraper"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/articles", func(w http.ResponseWriter, r *http.Request) {
		urlsParam := r.URL.Query().Get("urls")
		if urlsParam == "" {
			http.Error(w, "`urls` parameter is required", http.StatusBadRequest)
			return
		}

		urls := strings.Split(urlsParam, ",")

		articles, err := scraper.FetchArticles(urls)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(articles)
	})

	http.ListenAndServe(":8080", r)
}

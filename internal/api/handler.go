package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"tenome/internal/cerrors"
	"tenome/internal/service"
)

type Handler struct {
	searchService service.Searcher
	crawlService  service.Crawler
}

type CrawlRequest struct {
	URL string `json:"url"`
}

func NewHandler(searchService service.Searcher, crawlService service.Crawler) *Handler {
	return &Handler{searchService, crawlService}
}

func (h *Handler) Crawl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var crawlurl CrawlRequest

	if err := json.NewDecoder(r.Body).Decode(&crawlurl); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if crawlurl.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(crawlurl.URL); err != nil {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	err := h.crawlService.Submit(r.Context(), crawlurl.URL)

	if errors.Is(err, cerrors.ErrQueueFull) {
		http.Error(w, "job queue full", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "job submitted"})

}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("q")

	if term == "" {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}
	pages, err := h.searchService.Search(r.Context(), term)

	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pages)

}

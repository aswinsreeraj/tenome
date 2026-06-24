package api

import (
	"encoding/json"
	"net/http"
	"tenome/internal/service"
	"tenome/internal/worker"
)

type Handler struct {
	searchService *service.SearchService
	jobs          chan<- worker.CrawlJob
}

type CrawlRequest struct {
	URL string `json:"url"`
}

func (h *Handler) Crawl(w http.ResponseWriter, r *http.Request) {
	var crawlurl CrawlRequest
	if err := json.NewDecoder(r.Body).Decode(&crawlurl); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if crawlurl.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	select {
	case h.jobs <- worker.CrawlJob{URL: crawlurl.URL}:
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "job submitted",
		})
	default:
		http.Error(w, "job queue is full", http.StatusServiceUnavailable)
	}

}

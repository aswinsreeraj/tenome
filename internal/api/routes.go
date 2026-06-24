package api

import "net/http"

func NewRouter(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /crawl",
		handler.Crawl,
	)

	mux.HandleFunc(
		"GET /search",
		handler.Search,
	)

	return mux
}

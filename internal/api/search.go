package api

import (
	"encoding/json"
	"net/http"
)

type SearchRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

type SearchResponse struct {
	Query   string   `json:"query"`
	Results []string `json:"results"`
}

func Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request SearchRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if request.Query == "" {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	if request.Limit <= 0 {
		request.Limit = 10
	}

	response := SearchResponse{
		Query:   request.Query,
		Results: []string{},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

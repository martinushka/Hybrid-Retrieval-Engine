package search

import (
	"strings"

	"github.com/martinushka/ios-rag/internal/product"
)

func scoreProduct(query string, p product.Product) float64 {
	query = strings.ToLower(strings.TrimSpace(query))

	if query == "" {
		return 0
	}

	title := strings.ToLower(p.Title)
	description := strings.ToLower(p.Description)

	var score float64

	if title == query {
		score += 10
	}

	if strings.Contains(title, query) {
		score += 5
	}

	if strings.Contains(description, query) {
		score += 2
	}

	return score
}

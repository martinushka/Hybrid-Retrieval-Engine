package search

import (
	"github.com/martinushka/ios-rag/internal/product"
	"github.com/martinushka/ios-rag/internal/text"
)

func scoreProduct(query string, p product.Product) float64 {
	normalizedQuery := text.Normalize(query)
	normalizedTitle := text.Normalize(p.Title)

	queryTokens := text.Tokens(query)
	titleTokens := text.Tokens(p.Title)
	descriptionTokens := text.Tokens(p.Description)
	categoryTokens := text.Tokens(p.Category)

	if len(queryTokens) == 0 {
		return 0
	}

	var score float64

	if normalizedQuery == normalizedTitle {
		score += 12
	}

	titleSet := make(map[string]struct{}, len(titleTokens))
	descriptionSet := make(map[string]struct{}, len(descriptionTokens))
	categorySet := make(map[string]struct{}, len(categoryTokens))

	for _, token := range titleTokens {
		titleSet[token] = struct{}{}
	}

	for _, token := range descriptionTokens {
		descriptionSet[token] = struct{}{}
	}

	for _, token := range categoryTokens {
		categorySet[token] = struct{}{}
	}

	for _, token := range queryTokens {
		if _, ok := titleSet[token]; ok {
			score += 3
			continue
		}

		if _, ok := categorySet[token]; ok {
			score += 2
			continue
		}

		if _, ok := descriptionSet[token]; ok {
			score++
		}
	}

	return score
}

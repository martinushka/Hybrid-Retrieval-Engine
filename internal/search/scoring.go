package search

import (
	"github.com/martinushka/ios-rag/internal/product"
	"github.com/martinushka/ios-rag/internal/text"
)

const (
	maxLexicalScore   = 5.0
	semanticThreshold = 0.80
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

func normalizeLexicalScore(score float64) float64 {
	if score <= 0 {
		return 0
	}

	if score >= maxLexicalScore {
		return 1
	}

	return score / maxLexicalScore
}

func normalizeSemanticScore(score float64) float64 {
	if score <= semanticThreshold {
		return 0
	}

	return (score - semanticThreshold) / (1 - semanticThreshold)
}

func hybridScore(lexicalScore, semanticScore float64) float64 {
	const lexicalWeight = 0.6
	const semanticWeight = 0.4

	lexical := normalizeLexicalScore(lexicalScore)
	semantic := normalizeSemanticScore(semanticScore)

	return lexicalWeight*lexical + semanticWeight*semantic
}

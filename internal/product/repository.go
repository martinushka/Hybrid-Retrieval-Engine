package product

import (
	"context"
	"strings"

	"github.com/martinushka/ios-rag/internal/text"
)

type Candidate struct {
	Product      Product
	LexicalScore float64
}

type Repository interface {
	List(ctx context.Context) ([]Product, error)
	Search(ctx context.Context, query string, limit int) ([]Candidate, error)
}

type InMemoryRepository struct {
	products []Product
}

func NewInMemoryRepository(products []Product) *InMemoryRepository {
	return &InMemoryRepository{
		products: products,
	}
}

func (r *InMemoryRepository) List(ctx context.Context) ([]Product, error) {
	products := make([]Product, len(r.products))
	copy(products, r.products)

	return products, nil
}

func (r *InMemoryRepository) Search(
	ctx context.Context,
	query string,
	limit int,
) ([]Candidate, error) {
	query = strings.TrimSpace(query)

	if query == "" || limit <= 0 {
		return []Candidate{}, nil
	}

	queryTokens := text.Tokens(query)

	results := make([]Candidate, 0, limit)

	for _, p := range r.products {
		titleTokens := text.Tokens(p.Title)
		descriptionTokens := text.Tokens(p.Description)
		categoryTokens := text.Tokens(p.Category)

		if matchesAny(queryTokens, titleTokens) ||
			matchesAny(queryTokens, descriptionTokens) ||
			matchesAny(queryTokens, categoryTokens) {
			results = append(results, Candidate{
				Product:      p,
				LexicalScore: 0,
			})
		}

		if len(results) >= limit {
			break
		}
	}

	return results, nil
}

func matchesAny(queryTokens, documentTokens []string) bool {
	documentSet := make(map[string]struct{}, len(documentTokens))

	for _, token := range documentTokens {
		documentSet[token] = struct{}{}
	}

	for _, token := range queryTokens {
		if _, ok := documentSet[token]; ok {
			return true
		}
	}

	return false
}

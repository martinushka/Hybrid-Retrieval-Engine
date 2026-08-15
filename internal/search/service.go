package search

import (
	"context"
	"sort"
	"strings"

	"github.com/martinushka/ios-rag/internal/product"
)

type SearchResult struct {
	Product product.Product
	Score   float64
}

type Service interface {
	Search(ctx context.Context, query string, limit int) ([]SearchResult, error)
}

type InMemoryService struct {
	repository product.Repository
}

func NewInMemoryService(repository product.Repository) *InMemoryService {
	return &InMemoryService{
		repository: repository,
	}
}

type scoredProduct struct {
	product product.Product
	score   float64
}

func (s *InMemoryService) Search(
	ctx context.Context,
	query string,
	limit int,
) ([]SearchResult, error) {
	query = strings.TrimSpace(query)

	if query == "" || limit <= 0 {
		return []SearchResult{}, nil
	}

	candidateLimit := limit * 5
	if candidateLimit < 50 {
		candidateLimit = 50
	}

	candidates, err := s.repository.Search(ctx, query, candidateLimit)
	if err != nil {
		return nil, err
	}

	scored := make([]scoredProduct, 0, len(candidates))

	for _, candidate := range candidates {
		baseScore := scoreProduct(query, candidate.Product)

		hybridScore := baseScore + candidate.LexicalScore

		if hybridScore > 0 {
			scored = append(scored, scoredProduct{
				product: candidate.Product,
				score:   hybridScore,
			})
		}
	}

	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	if len(scored) > limit {
		scored = scored[:limit]
	}

	results := make([]SearchResult, 0, len(scored))

	for _, item := range scored {
		results = append(results, SearchResult{
			Product: item.product,
			Score:   item.score,
		})
	}

	return results, nil
}

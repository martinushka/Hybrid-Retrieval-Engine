package search

import (
	"context"
	"sort"
	"strings"

	"github.com/martinushka/ios-rag/internal/embedding"
	"github.com/martinushka/ios-rag/internal/product"
)

const rrfK = 60.0

type SearchResult struct {
	Product product.Product
	Score   float64
}

type Service interface {
	Search(ctx context.Context, query string, limit int) ([]SearchResult, error)
}

type HybridService struct {
	repository product.Repository
	embedder   embedding.Provider
}

func NewHybridService(
	repository product.Repository,
	embedder embedding.Provider,
) *HybridService {
	return &HybridService{
		repository: repository,
		embedder:   embedder,
	}
}

type scoredProduct struct {
	product product.Product
	score   float64
}

func rrfScore(rank int) float64 {
	return 1.0 / (rrfK + float64(rank))
}

func (s *HybridService) Search(
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

	candidates := make(map[string]scoredProduct)

	lexicalCandidates, err := s.repository.Search(
		ctx,
		query,
		candidateLimit,
	)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(lexicalCandidates, func(i, j int) bool {
		left := scoreProduct(query, lexicalCandidates[i].Product)
		right := scoreProduct(query, lexicalCandidates[j].Product)

		return left > right
	})

	for rank, candidate := range lexicalCandidates {
		candidates[candidate.Product.ID] = scoredProduct{
			product: candidate.Product,
			score:   rrfScore(rank + 1),
		}
	}

	if s.embedder != nil {
		queryEmbedding, err := s.embedder.Embed(
			ctx,
			"query: "+query,
		)
		if err != nil {
			return nil, err
		}

		semanticCandidates, err := s.repository.SemanticSearch(
			ctx,
			queryEmbedding,
			candidateLimit,
		)
		if err != nil {
			return nil, err
		}

		for rank, candidate := range semanticCandidates {
			item, exists := candidates[candidate.Product.ID]

			if !exists {
				item = scoredProduct{
					product: candidate.Product,
					score:   0,
				}
			}

			item.score += rrfScore(rank + 1)
			candidates[candidate.Product.ID] = item
		}
	}

	scored := make([]scoredProduct, 0, len(candidates))

	for _, item := range candidates {
		if item.score > 0 {
			scored = append(scored, item)
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

// InMemoryService сохраняет API, которое используют существующие тесты.
// Реализация использует тот же hybrid search pipeline.
type InMemoryService struct {
	repository product.Repository
	embedder   embedding.Provider
}

func NewInMemoryService(
	repository product.Repository,
	embedder embedding.Provider,
) *InMemoryService {
	return &InMemoryService{
		repository: repository,
		embedder:   embedder,
	}
}

func (s *InMemoryService) Search(
	ctx context.Context,
	query string,
	limit int,
) ([]SearchResult, error) {
	service := &HybridService{
		repository: s.repository,
		embedder:   s.embedder,
	}

	return service.Search(ctx, query, limit)
}

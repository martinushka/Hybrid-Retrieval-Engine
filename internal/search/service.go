package search

import (
	"context"
	"sort"
	"strings"

	"github.com/martinushka/ios-rag/internal/embedding"
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

	candidates := make(map[string]scoredProduct)

	// Lexical search.
	lexicalCandidates, err := s.repository.Search(
		ctx,
		query,
		candidateLimit,
	)
	if err != nil {
		return nil, err
	}

	for _, candidate := range lexicalCandidates {
		baseScore := scoreProduct(query, candidate.Product)

		candidates[candidate.Product.ID] = scoredProduct{
			product: candidate.Product,
			score:   0.6 * baseScore,
		}
	}

	// Semantic search.
	//
	// Semantic candidates can be added even when there is no
	// lexical match. This allows the search to find products
	// that are relevant by meaning rather than exact words.
	if s.embedder != nil {
		queryEmbedding, err := s.embedder.Embed(ctx, query)
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

		for _, candidate := range semanticCandidates {
			if candidate.SemanticScore < semanticThreshold {
				continue
			}

			item, exists := candidates[candidate.Product.ID]

			if !exists {
				item = scoredProduct{
					product: candidate.Product,
					score:   0,
				}
			}

			item.score += 0.4 * candidate.SemanticScore
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

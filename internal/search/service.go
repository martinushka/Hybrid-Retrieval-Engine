package search

import (
	"context"
	"sort"
	"strings"

	"github.com/martinushka/ios-rag/internal/product"
)

type Service interface {
	Search(ctx context.Context, query string, limit int) ([]string, error)
}

type InMemoryService struct {
	products []product.Product
}

func NewInMemoryService(products []product.Product) *InMemoryService {
	return &InMemoryService{
		products: products,
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
) ([]string, error) {
	query = strings.TrimSpace(query)

	if query == "" || limit <= 0 {
		return []string{}, nil
	}

	scored := make([]scoredProduct, 0, len(s.products))

	for _, p := range s.products {
		score := scoreProduct(query, p)

		if score > 0 {
			scored = append(scored, scoredProduct{
				product: p,
				score:   score,
			})
		}
	}

	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	if len(scored) > limit {
		scored = scored[:limit]
	}

	results := make([]string, 0, len(scored))

	for _, item := range scored {
		results = append(results, item.product.Title)
	}

	return results, nil
}

package search

import (
	"context"
	"testing"

	"github.com/martinushka/ios-rag/internal/product"
)

type fakeEmbeddingProvider struct{}

func (f *fakeEmbeddingProvider) Embed(
	ctx context.Context,
	text string,
) ([]float32, error) {
	return make([]float32, 384), nil
}

func TestInMemoryServiceSearch(t *testing.T) {
	products := []product.Product{
		{
			ID:          "1",
			Title:       "Lenovo ThinkPad X1",
			Description: "Business laptop",
			Category:    "laptops",
			Price:       1299.99,
		},
		{
			ID:          "2",
			Title:       "Lenovo",
			Description: "Laptop",
			Category:    "laptops",
			Price:       999.99,
		},
		{
			ID:          "3",
			Title:       "MacBook Air M2",
			Description: "Apple laptop",
			Category:    "laptops",
			Price:       999.99,
		},
	}

	repository := product.NewInMemoryRepository(products)
	service := NewInMemoryService(
		repository,
		&fakeEmbeddingProvider{},
	)

	results, err := service.Search(
		context.Background(),
		"Lenovo",
		10,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if results[0].Product.Title != "Lenovo" {
		t.Fatalf("expected exact title first, got %s", results[0].Product.Title)
	}

	if results[1].Product.Title != "Lenovo ThinkPad X1" {
		t.Fatalf("expected ThinkPad second, got %s", results[1].Product.Title)
	}

	if results[0].Score <= results[1].Score {
		t.Fatalf(
			"expected exact title to have higher score: %f <= %f",
			results[0].Score,
			results[1].Score,
		)
	}
}

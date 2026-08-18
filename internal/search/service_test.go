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

type fakeSearchRepository struct {
	lexicalCandidates  []product.Candidate
	semanticCandidates []product.Candidate
}

func (r *fakeSearchRepository) List(
	ctx context.Context,
) ([]product.Product, error) {
	return nil, nil
}

func (r *fakeSearchRepository) Search(
	ctx context.Context,
	query string,
	limit int,
) ([]product.Candidate, error) {
	if limit > len(r.lexicalCandidates) {
		limit = len(r.lexicalCandidates)
	}

	return r.lexicalCandidates[:limit], nil
}

func (r *fakeSearchRepository) SemanticSearch(
	ctx context.Context,
	embedding []float32,
	limit int,
) ([]product.Candidate, error) {
	if limit > len(r.semanticCandidates) {
		limit = len(r.semanticCandidates)
	}

	return r.semanticCandidates[:limit], nil
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

func TestRRFCombinesLexicalAndSemanticResults(t *testing.T) {
	repository := &fakeSearchRepository{
		lexicalCandidates: []product.Candidate{
			{
				Product: product.Product{
					ID:    "1",
					Title: "Lenovo ThinkPad X1",
				},
			},
			{
				Product: product.Product{
					ID:    "2",
					Title: "MacBook Air M2",
				},
			},
		},
		semanticCandidates: []product.Candidate{
			{
				Product: product.Product{
					ID:    "2",
					Title: "MacBook Air M2",
				},
				SemanticScore: 0.95,
			},
			{
				Product: product.Product{
					ID:    "3",
					Title: "Samsung Galaxy S24",
				},
				SemanticScore: 0.90,
			},
		},
	}

	service := NewInMemoryService(
		repository,
		&fakeEmbeddingProvider{},
	)

	results, err := service.Search(
		context.Background(),
		"ноутбук",
		10,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	if results[0].Product.ID != "2" {
		t.Fatalf(
			"expected product 2 to rank first because it appears in both lists, got %s",
			results[0].Product.ID,
		)
	}

	if results[0].Score <= results[1].Score {
		t.Fatalf(
			"expected combined RRF score to be highest: %f <= %f",
			results[0].Score,
			results[1].Score,
		)
	}
}

func TestRRFKeepsSemanticMatchesRegardlessOfRawSimilarity(t *testing.T) {
	repository := &fakeSearchRepository{
		semanticCandidates: []product.Candidate{
			{
				Product: product.Product{
					ID:    "1",
					Title: "Lenovo ThinkPad X1",
				},
				SemanticScore: 0.79,
			},
			{
				Product: product.Product{
					ID:    "2",
					Title: "MacBook Air M2",
				},
				SemanticScore: 0.81,
			},
		},
	}

	service := NewInMemoryService(
		repository,
		&fakeEmbeddingProvider{},
	)

	results, err := service.Search(
		context.Background(),
		"ноутбук",
		10,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestRRFAllowsSemanticOnlyResults(t *testing.T) {
	repository := &fakeSearchRepository{
		semanticCandidates: []product.Candidate{
			{
				Product: product.Product{
					ID:    "3",
					Title: "Samsung Galaxy S24",
				},
				SemanticScore: 0.95,
			},
		},
	}

	service := NewInMemoryService(
		repository,
		&fakeEmbeddingProvider{},
	)

	results, err := service.Search(
		context.Background(),
		"телефон",
		10,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 semantic-only result, got %d", len(results))
	}

	if results[0].Product.ID != "3" {
		t.Fatalf(
			"expected semantic-only product 3, got %s",
			results[0].Product.ID,
		)
	}
}

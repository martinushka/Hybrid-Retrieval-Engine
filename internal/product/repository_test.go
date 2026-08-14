package product

import (
	"context"
	"testing"
)

func TestInMemoryRepositoryList(t *testing.T) {
	products := []Product{
		{
			ID:          "1",
			Title:       "Lenovo ThinkPad X1",
			Description: "Business laptop",
			Category:    "laptops",
			Price:       1299.99,
		},
		{
			ID:          "2",
			Title:       "MacBook Air M2",
			Description: "Lightweight laptop",
			Category:    "laptops",
			Price:       999.99,
		},
	}

	repository := NewInMemoryRepository(products)

	got, err := repository.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 products, got %d", len(got))
	}

	if got[0].ID != "1" {
		t.Fatalf("expected first product ID 1, got %s", got[0].ID)
	}

	if got[1].ID != "2" {
		t.Fatalf("expected second product ID 2, got %s", got[1].ID)
	}
}

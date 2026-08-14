package search

import (
	"context"
	"testing"
)

func TestInMemoryServiceSearch(t *testing.T) {
	service := NewInMemoryService()

	results, err := service.Search(
		context.Background(),
		"ноутбук Lenovo",
		10,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

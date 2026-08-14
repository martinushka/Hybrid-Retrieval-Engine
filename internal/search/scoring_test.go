package search

import (
	"testing"

	"github.com/martinushka/ios-rag/internal/product"
)

func TestScoreProduct(t *testing.T) {
	p := product.Product{
		Title:       "Lenovo ThinkPad X1",
		Description: "Business laptop",
	}

	score := scoreProduct("Lenovo", p)

	if score != 5 {
		t.Fatalf("expected score 5, got %f", score)
	}
}

func TestScoreProductExactTitle(t *testing.T) {
	p := product.Product{
		Title: "Lenovo",
	}

	score := scoreProduct("Lenovo", p)

	if score != 15 {
		t.Fatalf("expected score 15, got %f", score)
	}
}

func TestScoreProductDescription(t *testing.T) {
	p := product.Product{
		Title:       "ThinkPad X1",
		Description: "Lenovo business laptop",
	}

	score := scoreProduct("Lenovo", p)

	if score != 2 {
		t.Fatalf("expected score 2, got %f", score)
	}
}

func TestScoreProductNoMatch(t *testing.T) {
	p := product.Product{
		Title:       "MacBook Air",
		Description: "Apple laptop",
	}

	score := scoreProduct("Lenovo", p)

	if score != 0 {
		t.Fatalf("expected score 0, got %f", score)
	}
}

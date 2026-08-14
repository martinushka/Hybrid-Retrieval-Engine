package product

import "testing"

func TestProduct(t *testing.T) {
	p := Product{
		ID:          "1",
		Title:       "Lenovo ThinkPad X1",
		Description: "14-inch business laptop",
		Category:    "laptops",
		Price:       1299.99,
	}

	if p.ID != "1" {
		t.Fatalf("expected ID 1, got %s", p.ID)
	}

	if p.Title != "Lenovo ThinkPad X1" {
		t.Fatalf("unexpected title: %s", p.Title)
	}

	if p.Category != "laptops" {
		t.Fatalf("unexpected category: %s", p.Category)
	}

	if p.Price != 1299.99 {
		t.Fatalf("unexpected price: %f", p.Price)
	}
}

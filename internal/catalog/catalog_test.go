package catalog

import "testing"

func TestSeed(t *testing.T) {
	products := Seed()

	if len(products) != 3 {
		t.Fatalf("expected 3 products, got %d", len(products))
	}

	if products[0].Title != "Lenovo ThinkPad X1" {
		t.Fatalf("unexpected first product: %s", products[0].Title)
	}
}
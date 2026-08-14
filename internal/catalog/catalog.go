package catalog

import "github.com/martinushka/ios-rag/internal/product"

func Seed() []product.Product {
	return []product.Product{
		{
			ID:          "1",
			Title:       "Lenovo ThinkPad X1",
			Description: "Lightweight business laptop for professionals",
			Category:    "laptops",
			Price:       1299.99,
		},
		{
			ID:          "2",
			Title:       "MacBook Air M2",
			Description: "Thin and lightweight laptop for everyday work",
			Category:    "laptops",
			Price:       999.99,
		},
		{
			ID:          "3",
			Title:       "ASUS Zenbook 14",
			Description: "Compact laptop with OLED display",
			Category:    "laptops",
			Price:       1099.99,
		},
	}
}
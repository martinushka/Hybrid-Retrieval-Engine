package product

import "context"

type Repository interface {
	List(ctx context.Context) ([]Product, error)
}

type InMemoryRepository struct {
	products []Product
}

func NewInMemoryRepository(products []Product) *InMemoryRepository {
	return &InMemoryRepository{
		products: products,
	}
}

func (r *InMemoryRepository) List(ctx context.Context) ([]Product, error) {
	products := make([]Product, len(r.products))
	copy(products, r.products)

	return products, nil
}

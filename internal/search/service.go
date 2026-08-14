package search

import "context"

type Service interface {
	Search(ctx context.Context, query string, limit int) ([]string, error)
}

type InMemoryService struct{}

func NewInMemoryService() *InMemoryService {
	return &InMemoryService{}
}

func (s *InMemoryService) Search(
	ctx context.Context,
	query string,
	limit int,
) ([]string, error) {
	return []string{}, nil
}

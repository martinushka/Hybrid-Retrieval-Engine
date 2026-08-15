package product

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	conn *pgx.Conn
}

func NewPostgresRepository(conn *pgx.Conn) *PostgresRepository {
	return &PostgresRepository{
		conn: conn,
	}
}

func (r *PostgresRepository) List(ctx context.Context) ([]Product, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT id, title, description, category, price
		FROM products
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)

	for rows.Next() {
		var p Product

		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.Category,
			&p.Price,
		); err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *PostgresRepository) Search(
	ctx context.Context,
	query string,
	limit int,
) ([]Candidate, error) {
	query = strings.TrimSpace(query)

	if query == "" || limit <= 0 {
		return []Candidate{}, nil
	}

	rows, err := r.conn.Query(ctx, `
		SELECT
			id,
			title,
			description,
			category,
			price,
			ts_rank(
				search_vector,
				websearch_to_tsquery('simple', $1)
			) AS lexical_score
		FROM products
		WHERE search_vector @@ websearch_to_tsquery('simple', $1)
		ORDER BY lexical_score DESC
		LIMIT $2
	`, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search products: %w", err)
	}
	defer rows.Close()

	results := make([]Candidate, 0, limit)

	for rows.Next() {
		var candidate Candidate

		if err := rows.Scan(
			&candidate.Product.ID,
			&candidate.Product.Title,
			&candidate.Product.Description,
			&candidate.Product.Category,
			&candidate.Product.Price,
			&candidate.LexicalScore,
		); err != nil {
			return nil, err
		}

		results = append(results, candidate)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

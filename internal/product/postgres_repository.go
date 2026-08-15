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
) ([]Product, error) {
	query = strings.TrimSpace(query)

	if query == "" || limit <= 0 {
		return []Product{}, nil
	}

	pattern := "%" + query + "%"

	rows, err := r.conn.Query(ctx, `
		SELECT id, title, description, category, price
		FROM products
		WHERE
			title ILIKE $1
			OR description ILIKE $1
			OR category ILIKE $1
		ORDER BY id
		LIMIT $2
	`, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("search products: %w", err)
	}
	defer rows.Close()

	products := make([]Product, 0, limit)

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

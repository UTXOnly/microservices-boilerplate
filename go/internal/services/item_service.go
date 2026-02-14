package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/microservices-boilerplate/go/internal/models"
)

// ErrNotFound is returned when a requested entity does not exist.
var ErrNotFound = errors.New("not found")

// ItemService provides async CRUD operations for items.
// Uses pgxpool for concurrent, non-blocking database access.
type ItemService struct {
	pool *pgxpool.Pool
}

// NewItemService creates a new ItemService.
func NewItemService(pool *pgxpool.Pool) *ItemService {
	return &ItemService{pool: pool}
}

// List returns items with pagination.
func (s *ItemService) List(ctx context.Context, skip, limit int) ([]*models.Item, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, description, created_at, updated_at
		FROM items
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, rows.Err()
}

// GetByID returns a single item by ID.
// Returns ErrNotFound if the item does not exist.
func (s *ItemService) GetByID(ctx context.Context, id int) (*models.Item, error) {
	var item models.Item
	err := s.pool.QueryRow(ctx, `
		SELECT id, name, description, created_at, updated_at
		FROM items WHERE id = $1
	`, id).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

// Create inserts a new item.
func (s *ItemService) Create(ctx context.Context, data *models.ItemCreate) (*models.Item, error) {
	var item models.Item
	err := s.pool.QueryRow(ctx, `
		INSERT INTO items (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`, data.Name, data.Description).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update performs a partial update (only non-nil fields are updated).
// Returns ErrNotFound if the item does not exist.
func (s *ItemService) Update(ctx context.Context, id int, data *models.ItemUpdate) (*models.Item, error) {
	// Fetch existing first for merge
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	name := existing.Name
	if data.Name != nil && *data.Name != "" {
		name = *data.Name
	}
	desc := existing.Description
	if data.Description != nil {
		desc = data.Description
	}

	var item models.Item
	err = s.pool.QueryRow(ctx, `
		UPDATE items SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, created_at, updated_at
	`, id, name, desc).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

// Delete removes an item.
func (s *ItemService) Delete(ctx context.Context, id int) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM items WHERE id = $1`, id)
	return err
}

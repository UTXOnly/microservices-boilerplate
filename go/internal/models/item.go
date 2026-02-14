package models

import "time"

// Item represents an item entity - replace or extend for your domain.
type Item struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ItemCreate is the payload for creating an item.
type ItemCreate struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// ItemUpdate is the payload for partial updates.
type ItemUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

package model

import "time"

type Deck struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	ParentID     *int64    `json:"parent_id,omitempty" db:"parent_id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	SRSAlgorithm *string   `json:"srs_algorithm,omitempty" db:"srs_algorithm"`
	Position     int       `json:"position" db:"position"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type DeckWithChildren struct {
	Deck
	Children []*DeckWithChildren `json:"children,omitempty"`
}

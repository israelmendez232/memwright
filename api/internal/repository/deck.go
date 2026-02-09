package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"memwright/api/internal/model"
)

type DeckRepository interface {
	Create(ctx context.Context, deck *model.Deck) error
	GetByID(ctx context.Context, id int64) (*model.Deck, error)
	GetByUserID(ctx context.Context, userID int64) ([]*model.Deck, error)
	Update(ctx context.Context, deck *model.Deck) error
	UpdateSRSConfig(ctx context.Context, id int64, config *model.SRSConfig) error
	Delete(ctx context.Context, id int64) error
}

type deckRepository struct {
	db DB
}

func NewDeckRepository(db DB) DeckRepository {
	return &deckRepository{db: db}
}

func (r *deckRepository) Create(ctx context.Context, deck *model.Deck) error {
	var configJSON []byte
	var err error
	if deck.SRSConfig != nil {
		configJSON, err = json.Marshal(deck.SRSConfig)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO decks (user_id, parent_id, name, description, algorithm, srs_config, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		deck.UserID,
		deck.ParentID,
		deck.Name,
		deck.Description,
		deck.Algorithm,
		configJSON,
		deck.Position,
	).Scan(&deck.ID, &deck.CreatedAt, &deck.UpdatedAt)
}

func (r *deckRepository) GetByID(ctx context.Context, id int64) (*model.Deck, error) {
	query := `
		SELECT id, user_id, parent_id, name, description, algorithm, srs_config, position, created_at, updated_at
		FROM decks
		WHERE id = $1`

	deck := &model.Deck{}
	var configJSON sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&deck.ID,
		&deck.UserID,
		&deck.ParentID,
		&deck.Name,
		&deck.Description,
		&deck.Algorithm,
		&configJSON,
		&deck.Position,
		&deck.CreatedAt,
		&deck.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if configJSON.Valid && configJSON.String != "" {
		deck.SRSConfig = &model.SRSConfig{}
		if err := json.Unmarshal([]byte(configJSON.String), deck.SRSConfig); err != nil {
			return nil, err
		}
	}

	return deck, nil
}

func (r *deckRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.Deck, error) {
	query := `
		SELECT id, user_id, parent_id, name, description, algorithm, srs_config, position, created_at, updated_at
		FROM decks
		WHERE user_id = $1
		ORDER BY position, name`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decks []*model.Deck
	for rows.Next() {
		deck := &model.Deck{}
		var configJSON sql.NullString

		err := rows.Scan(
			&deck.ID,
			&deck.UserID,
			&deck.ParentID,
			&deck.Name,
			&deck.Description,
			&deck.Algorithm,
			&configJSON,
			&deck.Position,
			&deck.CreatedAt,
			&deck.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if configJSON.Valid && configJSON.String != "" {
			deck.SRSConfig = &model.SRSConfig{}
			if err := json.Unmarshal([]byte(configJSON.String), deck.SRSConfig); err != nil {
				return nil, err
			}
		}

		decks = append(decks, deck)
	}

	return decks, rows.Err()
}

func (r *deckRepository) Update(ctx context.Context, deck *model.Deck) error {
	var configJSON []byte
	var err error
	if deck.SRSConfig != nil {
		configJSON, err = json.Marshal(deck.SRSConfig)
		if err != nil {
			return err
		}
	}

	query := `
		UPDATE decks
		SET parent_id = $2, name = $3, description = $4, algorithm = $5, srs_config = $6, position = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err = r.db.QueryRowContext(ctx, query,
		deck.ID,
		deck.ParentID,
		deck.Name,
		deck.Description,
		deck.Algorithm,
		configJSON,
		deck.Position,
	).Scan(&deck.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	return err
}

func (r *deckRepository) UpdateSRSConfig(ctx context.Context, id int64, config *model.SRSConfig) error {
	var configJSON []byte
	var err error
	if config != nil {
		configJSON, err = json.Marshal(config)
		if err != nil {
			return err
		}
	}

	query := `
		UPDATE decks
		SET srs_config = $2, updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, configJSON)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	return nil
}

func (r *deckRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM decks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	return nil
}

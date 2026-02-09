package repository

import (
	"context"
	"database/sql"
	"errors"

	"memwright/api/internal/model"
)

type CardRepository interface {
	Create(ctx context.Context, card *model.Card) error
	GetByID(ctx context.Context, id int64) (*model.Card, error)
	GetByDeckID(ctx context.Context, deckID int64) ([]*model.Card, error)
	Update(ctx context.Context, card *model.Card) error
	Delete(ctx context.Context, id int64) error
}

type cardRepository struct {
	db DB
}

func NewCardRepository(db DB) CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) Create(ctx context.Context, card *model.Card) error {
	query := `
		INSERT INTO cards (deck_id, type, front, back, extra, tags, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		card.DeckID,
		card.Type,
		card.Front,
		card.Back,
		card.Extra,
		tagsToArray(card.Tags),
		card.Position,
	).Scan(&card.ID, &card.CreatedAt, &card.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *cardRepository) GetByID(ctx context.Context, id int64) (*model.Card, error) {
	query := `
		SELECT id, deck_id, type, front, back, extra, tags, position, created_at, updated_at
		FROM cards
		WHERE id = $1`

	card := &model.Card{}
	var tags sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&card.ID,
		&card.DeckID,
		&card.Type,
		&card.Front,
		&card.Back,
		&card.Extra,
		&tags,
		&card.Position,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	card.Tags = arrayToTags(tags)
	return card, nil
}

func (r *cardRepository) GetByDeckID(ctx context.Context, deckID int64) ([]*model.Card, error) {
	query := `
		SELECT id, deck_id, type, front, back, extra, tags, position, created_at, updated_at
		FROM cards
		WHERE deck_id = $1
		ORDER BY position, id`

	rows, err := r.db.QueryContext(ctx, query, deckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*model.Card
	for rows.Next() {
		card := &model.Card{}
		var tags sql.NullString
		err := rows.Scan(
			&card.ID,
			&card.DeckID,
			&card.Type,
			&card.Front,
			&card.Back,
			&card.Extra,
			&tags,
			&card.Position,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		card.Tags = arrayToTags(tags)
		cards = append(cards, card)
	}
	return cards, rows.Err()
}

func (r *cardRepository) Update(ctx context.Context, card *model.Card) error {
	query := `
		UPDATE cards
		SET deck_id = $2, type = $3, front = $4, back = $5, extra = $6, tags = $7, position = $8, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		card.ID,
		card.DeckID,
		card.Type,
		card.Front,
		card.Back,
		card.Extra,
		tagsToArray(card.Tags),
		card.Position,
	).Scan(&card.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *cardRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM cards WHERE id = $1`

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

func tagsToArray(tags []string) interface{} {
	if len(tags) == 0 {
		return nil
	}
	return tags
}

func arrayToTags(tags sql.NullString) []string {
	if !tags.Valid || tags.String == "" || tags.String == "{}" {
		return nil
	}
	return parsePostgresArray(tags.String)
}

func parsePostgresArray(s string) []string {
	if len(s) < 2 {
		return nil
	}
	s = s[1 : len(s)-1]
	if s == "" {
		return nil
	}

	var result []string
	var current string
	inQuotes := false

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '"' {
			inQuotes = !inQuotes
		} else if ch == ',' && !inQuotes {
			result = append(result, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

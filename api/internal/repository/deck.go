package repository

import (
	"context"
	"database/sql"
	"errors"

	"memwright/api/internal/model"
)

type DeckRepository interface {
	Create(ctx context.Context, deck *model.Deck) error
	GetByID(ctx context.Context, id int64) (*model.Deck, error)
	GetByUserID(ctx context.Context, userID int64) ([]*model.Deck, error)
	GetSubdecks(ctx context.Context, parentID int64) ([]*model.Deck, error)
	GetAllSubdecks(ctx context.Context, parentID int64) ([]*model.Deck, error)
	GetDeckTree(ctx context.Context, userID int64) ([]*model.DeckWithChildren, error)
	Update(ctx context.Context, deck *model.Deck) error
	Delete(ctx context.Context, id int64) error
}

type deckRepository struct {
	db DB
}

func NewDeckRepository(db DB) DeckRepository {
	return &deckRepository{db: db}
}

func (r *deckRepository) Create(ctx context.Context, deck *model.Deck) error {
	query := `
		INSERT INTO decks (user_id, parent_id, name, description, srs_algorithm, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		deck.UserID,
		deck.ParentID,
		deck.Name,
		deck.Description,
		deck.SRSAlgorithm,
		deck.Position,
	).Scan(&deck.ID, &deck.CreatedAt, &deck.UpdatedAt)

	if err != nil {
		if isDuplicateKeyError(err, "decks_user_id_parent_id_name_key") {
			return model.ErrDuplicateName
		}
		return err
	}
	return nil
}

func (r *deckRepository) GetByID(ctx context.Context, id int64) (*model.Deck, error) {
	query := `
		SELECT id, user_id, parent_id, name, description, srs_algorithm, position, created_at, updated_at
		FROM decks
		WHERE id = $1`

	deck := &model.Deck{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&deck.ID,
		&deck.UserID,
		&deck.ParentID,
		&deck.Name,
		&deck.Description,
		&deck.SRSAlgorithm,
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
	return deck, nil
}

func (r *deckRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.Deck, error) {
	query := `
		SELECT id, user_id, parent_id, name, description, srs_algorithm, position, created_at, updated_at
		FROM decks
		WHERE user_id = $1
		ORDER BY position, name`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanDecks(rows)
}

func (r *deckRepository) GetSubdecks(ctx context.Context, parentID int64) ([]*model.Deck, error) {
	query := `
		SELECT id, user_id, parent_id, name, description, srs_algorithm, position, created_at, updated_at
		FROM decks
		WHERE parent_id = $1
		ORDER BY position, name`

	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanDecks(rows)
}

func (r *deckRepository) GetAllSubdecks(ctx context.Context, parentID int64) ([]*model.Deck, error) {
	query := `
		WITH RECURSIVE subdeck_tree AS (
			SELECT id, user_id, parent_id, name, description, srs_algorithm, position, created_at, updated_at
			FROM decks
			WHERE parent_id = $1

			UNION ALL

			SELECT d.id, d.user_id, d.parent_id, d.name, d.description, d.srs_algorithm, d.position, d.created_at, d.updated_at
			FROM decks d
			INNER JOIN subdeck_tree st ON d.parent_id = st.id
		)
		SELECT * FROM subdeck_tree
		ORDER BY parent_id NULLS FIRST, position, name`

	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanDecks(rows)
}

func (r *deckRepository) GetDeckTree(ctx context.Context, userID int64) ([]*model.DeckWithChildren, error) {
	decks, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return buildDeckTree(decks), nil
}

func (r *deckRepository) Update(ctx context.Context, deck *model.Deck) error {
	query := `
		UPDATE decks
		SET parent_id = $2, name = $3, description = $4, srs_algorithm = $5, position = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		deck.ID,
		deck.ParentID,
		deck.Name,
		deck.Description,
		deck.SRSAlgorithm,
		deck.Position,
	).Scan(&deck.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	if err != nil {
		if isDuplicateKeyError(err, "decks_user_id_parent_id_name_key") {
			return model.ErrDuplicateName
		}
		return err
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

func scanDecks(rows *sql.Rows) ([]*model.Deck, error) {
	var decks []*model.Deck
	for rows.Next() {
		deck := &model.Deck{}
		err := rows.Scan(
			&deck.ID,
			&deck.UserID,
			&deck.ParentID,
			&deck.Name,
			&deck.Description,
			&deck.SRSAlgorithm,
			&deck.Position,
			&deck.CreatedAt,
			&deck.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		decks = append(decks, deck)
	}
	return decks, rows.Err()
}

func buildDeckTree(decks []*model.Deck) []*model.DeckWithChildren {
	deckMap := make(map[int64]*model.DeckWithChildren)
	var roots []*model.DeckWithChildren

	for _, d := range decks {
		deckMap[d.ID] = &model.DeckWithChildren{Deck: *d}
	}

	for _, d := range decks {
		dwc := deckMap[d.ID]
		if d.ParentID == nil {
			roots = append(roots, dwc)
		} else if parent, ok := deckMap[*d.ParentID]; ok {
			parent.Children = append(parent.Children, dwc)
		}
	}

	return roots
}

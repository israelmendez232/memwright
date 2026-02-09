package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"memwright/api/internal/model"
)

type CardScheduleRepository interface {
	Create(ctx context.Context, schedule *model.CardSchedule) error
	GetByID(ctx context.Context, id int64) (*model.CardSchedule, error)
	GetByCardAndUser(ctx context.Context, cardID, userID int64) (*model.CardSchedule, error)
	GetDueCards(ctx context.Context, userID int64, deckID int64, dueBy time.Time, limit int) ([]*model.CardSchedule, error)
	GetNewCards(ctx context.Context, userID int64, deckID int64, limit int) ([]*model.CardSchedule, error)
	Update(ctx context.Context, schedule *model.CardSchedule) error
	Delete(ctx context.Context, id int64) error
}

type cardScheduleRepository struct {
	db DB
}

func NewCardScheduleRepository(db DB) CardScheduleRepository {
	return &cardScheduleRepository{db: db}
}

func (r *cardScheduleRepository) Create(ctx context.Context, schedule *model.CardSchedule) error {
	query := `
		INSERT INTO card_schedules (card_id, user_id, state, due_at, interval, ease_factor, review_count, lapse_count, last_reviewed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		schedule.CardID,
		schedule.UserID,
		schedule.State,
		schedule.DueAt,
		schedule.Interval,
		schedule.EaseFactor,
		schedule.ReviewCount,
		schedule.LapseCount,
		schedule.LastReviewedAt,
	).Scan(&schedule.ID, &schedule.CreatedAt, &schedule.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *cardScheduleRepository) GetByID(ctx context.Context, id int64) (*model.CardSchedule, error) {
	query := `
		SELECT id, card_id, user_id, state, due_at, interval, ease_factor, review_count, lapse_count, last_reviewed_at, created_at, updated_at
		FROM card_schedules
		WHERE id = $1`

	schedule := &model.CardSchedule{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&schedule.ID,
		&schedule.CardID,
		&schedule.UserID,
		&schedule.State,
		&schedule.DueAt,
		&schedule.Interval,
		&schedule.EaseFactor,
		&schedule.ReviewCount,
		&schedule.LapseCount,
		&schedule.LastReviewedAt,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func (r *cardScheduleRepository) GetByCardAndUser(ctx context.Context, cardID, userID int64) (*model.CardSchedule, error) {
	query := `
		SELECT id, card_id, user_id, state, due_at, interval, ease_factor, review_count, lapse_count, last_reviewed_at, created_at, updated_at
		FROM card_schedules
		WHERE card_id = $1 AND user_id = $2`

	schedule := &model.CardSchedule{}
	err := r.db.QueryRowContext(ctx, query, cardID, userID).Scan(
		&schedule.ID,
		&schedule.CardID,
		&schedule.UserID,
		&schedule.State,
		&schedule.DueAt,
		&schedule.Interval,
		&schedule.EaseFactor,
		&schedule.ReviewCount,
		&schedule.LapseCount,
		&schedule.LastReviewedAt,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func (r *cardScheduleRepository) GetDueCards(ctx context.Context, userID int64, deckID int64, dueBy time.Time, limit int) ([]*model.CardSchedule, error) {
	query := `
		SELECT cs.id, cs.card_id, cs.user_id, cs.state, cs.due_at, cs.interval, cs.ease_factor, cs.review_count, cs.lapse_count, cs.last_reviewed_at, cs.created_at, cs.updated_at
		FROM card_schedules cs
		INNER JOIN cards c ON cs.card_id = c.id
		WHERE cs.user_id = $1
			AND c.deck_id = $2
			AND cs.due_at <= $3
			AND cs.state IN ('learning', 'review', 'relearning')
		ORDER BY cs.due_at ASC
		LIMIT $4`

	rows, err := r.db.QueryContext(ctx, query, userID, deckID, dueBy, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCardSchedules(rows)
}

func (r *cardScheduleRepository) GetNewCards(ctx context.Context, userID int64, deckID int64, limit int) ([]*model.CardSchedule, error) {
	query := `
		SELECT cs.id, cs.card_id, cs.user_id, cs.state, cs.due_at, cs.interval, cs.ease_factor, cs.review_count, cs.lapse_count, cs.last_reviewed_at, cs.created_at, cs.updated_at
		FROM card_schedules cs
		INNER JOIN cards c ON cs.card_id = c.id
		WHERE cs.user_id = $1
			AND c.deck_id = $2
			AND cs.state = 'new'
		ORDER BY c.position, cs.id
		LIMIT $3`

	rows, err := r.db.QueryContext(ctx, query, userID, deckID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCardSchedules(rows)
}

func (r *cardScheduleRepository) Update(ctx context.Context, schedule *model.CardSchedule) error {
	query := `
		UPDATE card_schedules
		SET state = $2, due_at = $3, interval = $4, ease_factor = $5, review_count = $6, lapse_count = $7, last_reviewed_at = $8, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		schedule.ID,
		schedule.State,
		schedule.DueAt,
		schedule.Interval,
		schedule.EaseFactor,
		schedule.ReviewCount,
		schedule.LapseCount,
		schedule.LastReviewedAt,
	).Scan(&schedule.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *cardScheduleRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM card_schedules WHERE id = $1`

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

func scanCardSchedules(rows *sql.Rows) ([]*model.CardSchedule, error) {
	var schedules []*model.CardSchedule
	for rows.Next() {
		schedule := &model.CardSchedule{}
		err := rows.Scan(
			&schedule.ID,
			&schedule.CardID,
			&schedule.UserID,
			&schedule.State,
			&schedule.DueAt,
			&schedule.Interval,
			&schedule.EaseFactor,
			&schedule.ReviewCount,
			&schedule.LapseCount,
			&schedule.LastReviewedAt,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, rows.Err()
}

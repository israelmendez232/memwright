package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"memwright/api/internal/model"
)

type ReviewLogRepository interface {
	Create(ctx context.Context, log *model.ReviewLog) error
	GetByID(ctx context.Context, id int64) (*model.ReviewLog, error)
	GetByUserID(ctx context.Context, userID int64, limit int) ([]*model.ReviewLog, error)
	GetByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]*model.ReviewLog, error)
}

type reviewLogRepository struct {
	db DB
}

func NewReviewLogRepository(db DB) ReviewLogRepository {
	return &reviewLogRepository{db: db}
}

func (r *reviewLogRepository) Create(ctx context.Context, log *model.ReviewLog) error {
	query := `
		INSERT INTO review_logs (card_schedule_id, user_id, rating, previous_state, new_state, previous_ease, new_ease, previous_interval, new_interval, review_duration, reviewed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		log.CardScheduleID,
		log.UserID,
		log.Rating,
		log.PreviousState,
		log.NewState,
		log.PreviousEase,
		log.NewEase,
		log.PreviousInterval,
		log.NewInterval,
		log.ReviewDuration,
		log.ReviewedAt,
	).Scan(&log.ID)

	if err != nil {
		return err
	}
	return nil
}

func (r *reviewLogRepository) GetByID(ctx context.Context, id int64) (*model.ReviewLog, error) {
	query := `
		SELECT id, card_schedule_id, user_id, rating, previous_state, new_state, previous_ease, new_ease, previous_interval, new_interval, review_duration, reviewed_at
		FROM review_logs
		WHERE id = $1`

	log := &model.ReviewLog{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID,
		&log.CardScheduleID,
		&log.UserID,
		&log.Rating,
		&log.PreviousState,
		&log.NewState,
		&log.PreviousEase,
		&log.NewEase,
		&log.PreviousInterval,
		&log.NewInterval,
		&log.ReviewDuration,
		&log.ReviewedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return log, nil
}

func (r *reviewLogRepository) GetByUserID(ctx context.Context, userID int64, limit int) ([]*model.ReviewLog, error) {
	query := `
		SELECT id, card_schedule_id, user_id, rating, previous_state, new_state, previous_ease, new_ease, previous_interval, new_interval, review_duration, reviewed_at
		FROM review_logs
		WHERE user_id = $1
		ORDER BY reviewed_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReviewLogs(rows)
}

func (r *reviewLogRepository) GetByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]*model.ReviewLog, error) {
	query := `
		SELECT id, card_schedule_id, user_id, rating, previous_state, new_state, previous_ease, new_ease, previous_interval, new_interval, review_duration, reviewed_at
		FROM review_logs
		WHERE user_id = $1 AND reviewed_at >= $2 AND reviewed_at < $3
		ORDER BY reviewed_at ASC`

	rows, err := r.db.QueryContext(ctx, query, userID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReviewLogs(rows)
}

func scanReviewLogs(rows *sql.Rows) ([]*model.ReviewLog, error) {
	var logs []*model.ReviewLog
	for rows.Next() {
		log := &model.ReviewLog{}
		err := rows.Scan(
			&log.ID,
			&log.CardScheduleID,
			&log.UserID,
			&log.Rating,
			&log.PreviousState,
			&log.NewState,
			&log.PreviousEase,
			&log.NewEase,
			&log.PreviousInterval,
			&log.NewInterval,
			&log.ReviewDuration,
			&log.ReviewedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}

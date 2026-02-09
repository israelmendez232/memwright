package repository

import (
	"context"
	"database/sql"
	"errors"

	"memwright/api/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db DB
}

func NewUserRepository(db DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (email, password_hash, display_name, srs_algorithm, daily_new_cards, daily_reviews, timezone_offset, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		user.SRSAlgorithm,
		user.DailyNewCards,
		user.DailyReviews,
		user.TimezoneOffset,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if isDuplicateKeyError(err, "users_email_key") {
			return model.ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, display_name, srs_algorithm, daily_new_cards, daily_reviews, timezone_offset, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.SRSAlgorithm,
		&user.DailyNewCards,
		&user.DailyReviews,
		&user.TimezoneOffset,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, display_name, srs_algorithm, daily_new_cards, daily_reviews, timezone_offset, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.SRSAlgorithm,
		&user.DailyNewCards,
		&user.DailyReviews,
		&user.TimezoneOffset,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET email = $2, display_name = $3, srs_algorithm = $4, daily_new_cards = $5, daily_reviews = $6, timezone_offset = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		user.ID,
		user.Email,
		user.DisplayName,
		user.SRSAlgorithm,
		user.DailyNewCards,
		user.DailyReviews,
		user.TimezoneOffset,
	).Scan(&user.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	if err != nil {
		if isDuplicateKeyError(err, "users_email_key") {
			return model.ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

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

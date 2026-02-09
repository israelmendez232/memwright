package model

import "time"

type User struct {
	ID             int64     `json:"id" db:"id"`
	Email          string    `json:"email" db:"email"`
	PasswordHash   string    `json:"-" db:"password_hash"`
	DisplayName    string    `json:"display_name" db:"display_name"`
	SRSAlgorithm   string    `json:"srs_algorithm" db:"srs_algorithm"`
	DailyNewCards  int       `json:"daily_new_cards" db:"daily_new_cards"`
	DailyReviews   int       `json:"daily_reviews" db:"daily_reviews"`
	TimezoneOffset int       `json:"timezone_offset" db:"timezone_offset"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

const (
	SRSAlgorithmSM2  = "sm2"
	SRSAlgorithmFSRS = "fsrs"
)

const (
	DefaultDailyNewCards = 20
	DefaultDailyReviews  = 200
)

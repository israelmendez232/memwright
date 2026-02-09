package model

import "time"

type ScheduleState string

const (
	ScheduleStateNew        ScheduleState = "new"
	ScheduleStateLearning   ScheduleState = "learning"
	ScheduleStateReview     ScheduleState = "review"
	ScheduleStateRelearning ScheduleState = "relearning"
	ScheduleStateMastered   ScheduleState = "mastered"
)

type CardSchedule struct {
	ID             int64         `json:"id" db:"id"`
	CardID         int64         `json:"card_id" db:"card_id"`
	UserID         int64         `json:"user_id" db:"user_id"`
	State          ScheduleState `json:"state" db:"state"`
	DueAt          time.Time     `json:"due_at" db:"due_at"`
	Interval       int           `json:"interval" db:"interval"`
	EaseFactor     float64       `json:"ease_factor" db:"ease_factor"`
	ReviewCount    int           `json:"review_count" db:"review_count"`
	LapseCount     int           `json:"lapse_count" db:"lapse_count"`
	LastReviewedAt *time.Time    `json:"last_reviewed_at,omitempty" db:"last_reviewed_at"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
}

const DefaultEaseFactor = 2.5

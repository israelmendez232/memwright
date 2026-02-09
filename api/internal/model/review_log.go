package model

import "time"

type ReviewRating int

const (
	ReviewRatingWrong   ReviewRating = 1
	ReviewRatingCorrect ReviewRating = 2
	ReviewRatingEasy    ReviewRating = 3
)

type ReviewLog struct {
	ID              int64         `json:"id" db:"id"`
	CardScheduleID  int64         `json:"card_schedule_id" db:"card_schedule_id"`
	UserID          int64         `json:"user_id" db:"user_id"`
	Rating          ReviewRating  `json:"rating" db:"rating"`
	PreviousState   ScheduleState `json:"previous_state" db:"previous_state"`
	NewState        ScheduleState `json:"new_state" db:"new_state"`
	PreviousEase    float64       `json:"previous_ease" db:"previous_ease"`
	NewEase         float64       `json:"new_ease" db:"new_ease"`
	PreviousInterval int          `json:"previous_interval" db:"previous_interval"`
	NewInterval      int          `json:"new_interval" db:"new_interval"`
	ReviewDuration   int          `json:"review_duration" db:"review_duration"`
	ReviewedAt       time.Time    `json:"reviewed_at" db:"reviewed_at"`
}

type StudySession struct {
	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	DeckID         int64     `json:"deck_id" db:"deck_id"`
	CardsStudied   int       `json:"cards_studied" db:"cards_studied"`
	CardsCorrect   int       `json:"cards_correct" db:"cards_correct"`
	CardsWrong     int       `json:"cards_wrong" db:"cards_wrong"`
	TotalDuration  int       `json:"total_duration" db:"total_duration"`
	StartedAt      time.Time `json:"started_at" db:"started_at"`
	EndedAt        time.Time `json:"ended_at" db:"ended_at"`
}

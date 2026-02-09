package srs

import "time"

type Rating int

const (
	RatingWrong   Rating = 1
	RatingCorrect Rating = 2
	RatingEasy    Rating = 3
)

type State string

const (
	StateNew        State = "new"
	StateLearning   State = "learning"
	StateReview     State = "review"
	StateRelearning State = "relearning"
	StateMastered   State = "mastered"
)

type ScheduleInput struct {
	State          State
	Interval       int
	EaseFactor     float64
	ReviewCount    int
	LapseCount     int
	LastReviewedAt *time.Time
}

type ScheduleOutput struct {
	State      State
	Interval   int
	EaseFactor float64
	DueAt      time.Time
}

type Algorithm interface {
	Schedule(input ScheduleInput, rating Rating, now time.Time) ScheduleOutput
	Name() string
}

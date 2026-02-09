package srs

import (
	"math"
	"time"
)

type SM2Config struct {
	InitialEaseFactor  float64 `json:"initial_ease_factor" db:"initial_ease_factor"`
	MinEaseFactor      float64 `json:"min_ease_factor" db:"min_ease_factor"`
	MaxEaseFactor      float64 `json:"max_ease_factor" db:"max_ease_factor"`
	EaseDecrement      float64 `json:"ease_decrement" db:"ease_decrement"`
	EaseIncrement      float64 `json:"ease_increment" db:"ease_increment"`
	EasyBonusMultipler float64 `json:"easy_bonus_multiplier" db:"easy_bonus_multiplier"`
	GraduatingInterval int     `json:"graduating_interval" db:"graduating_interval"`
	MasteredThreshold  int     `json:"mastered_threshold" db:"mastered_threshold"`
}

type SM2 struct {
	config SM2Config
}

func NewSM2(config SM2Config) *SM2 {
	return &SM2{config: config}
}

func (s *SM2) Name() string {
	return "sm2"
}

func (s *SM2) Config() SM2Config {
	return s.config
}

func (s *SM2) Schedule(input ScheduleInput, rating Rating, now time.Time) ScheduleOutput {
	ease := input.EaseFactor
	if ease == 0 {
		ease = s.config.InitialEaseFactor
	}

	switch input.State {
	case StateNew:
		return s.scheduleNew(ease, rating, now)
	case StateLearning:
		return s.scheduleLearning(input, ease, rating, now)
	case StateReview, StateRelearning:
		return s.scheduleReview(input, ease, rating, now)
	case StateMastered:
		return s.scheduleMastered(input, ease, rating, now)
	default:
		return s.scheduleNew(ease, rating, now)
	}
}

func (s *SM2) scheduleNew(ease float64, rating Rating, now time.Time) ScheduleOutput {
	cfg := s.config
	switch rating {
	case RatingWrong:
		return ScheduleOutput{
			State:      StateNew,
			Interval:   0,
			EaseFactor: maxFloat(ease-cfg.EaseDecrement, cfg.MinEaseFactor),
			DueAt:      now,
		}
	case RatingCorrect:
		return ScheduleOutput{
			State:      StateLearning,
			Interval:   cfg.GraduatingInterval,
			EaseFactor: ease,
			DueAt:      now.AddDate(0, 0, cfg.GraduatingInterval),
		}
	case RatingEasy:
		interval := int(float64(cfg.GraduatingInterval) * cfg.EasyBonusMultipler * ease)
		return ScheduleOutput{
			State:      StateReview,
			Interval:   interval,
			EaseFactor: minFloat(ease+cfg.EaseIncrement, cfg.MaxEaseFactor),
			DueAt:      now.AddDate(0, 0, interval),
		}
	default:
		return ScheduleOutput{
			State:      StateNew,
			Interval:   0,
			EaseFactor: ease,
			DueAt:      now,
		}
	}
}

func (s *SM2) scheduleLearning(input ScheduleInput, ease float64, rating Rating, now time.Time) ScheduleOutput {
	cfg := s.config
	switch rating {
	case RatingWrong:
		return ScheduleOutput{
			State:      StateLearning,
			Interval:   cfg.GraduatingInterval,
			EaseFactor: maxFloat(ease-cfg.EaseDecrement, cfg.MinEaseFactor),
			DueAt:      now.AddDate(0, 0, cfg.GraduatingInterval),
		}
	case RatingCorrect:
		newInterval := maxInt(int(math.Round(float64(input.Interval)*ease)), input.Interval+1)
		return ScheduleOutput{
			State:      StateReview,
			Interval:   newInterval,
			EaseFactor: ease,
			DueAt:      now.AddDate(0, 0, newInterval),
		}
	case RatingEasy:
		newInterval := int(math.Round(float64(input.Interval) * ease * cfg.EasyBonusMultipler))
		newInterval = maxInt(newInterval, input.Interval+1)
		return ScheduleOutput{
			State:      StateReview,
			Interval:   newInterval,
			EaseFactor: minFloat(ease+cfg.EaseIncrement, cfg.MaxEaseFactor),
			DueAt:      now.AddDate(0, 0, newInterval),
		}
	default:
		return ScheduleOutput{
			State:      StateLearning,
			Interval:   input.Interval,
			EaseFactor: ease,
			DueAt:      now.AddDate(0, 0, input.Interval),
		}
	}
}

func (s *SM2) scheduleReview(input ScheduleInput, ease float64, rating Rating, now time.Time) ScheduleOutput {
	cfg := s.config
	switch rating {
	case RatingWrong:
		return ScheduleOutput{
			State:      StateRelearning,
			Interval:   cfg.GraduatingInterval,
			EaseFactor: maxFloat(ease-cfg.EaseDecrement, cfg.MinEaseFactor),
			DueAt:      now.AddDate(0, 0, cfg.GraduatingInterval),
		}
	case RatingCorrect:
		newInterval := int(math.Round(float64(input.Interval) * ease))
		newInterval = maxInt(newInterval, input.Interval+1)
		newState := StateReview
		if newInterval >= cfg.MasteredThreshold {
			newState = StateMastered
		}
		return ScheduleOutput{
			State:      newState,
			Interval:   newInterval,
			EaseFactor: ease,
			DueAt:      now.AddDate(0, 0, newInterval),
		}
	case RatingEasy:
		newInterval := int(math.Round(float64(input.Interval) * ease * cfg.EasyBonusMultipler))
		newInterval = maxInt(newInterval, input.Interval+1)
		newState := StateReview
		if newInterval >= cfg.MasteredThreshold {
			newState = StateMastered
		}
		return ScheduleOutput{
			State:      newState,
			Interval:   newInterval,
			EaseFactor: minFloat(ease+cfg.EaseIncrement, cfg.MaxEaseFactor),
			DueAt:      now.AddDate(0, 0, newInterval),
		}
	default:
		return ScheduleOutput{
			State:      StateReview,
			Interval:   input.Interval,
			EaseFactor: ease,
			DueAt:      now.AddDate(0, 0, input.Interval),
		}
	}
}

func (s *SM2) scheduleMastered(input ScheduleInput, ease float64, rating Rating, now time.Time) ScheduleOutput {
	cfg := s.config
	switch rating {
	case RatingWrong:
		return ScheduleOutput{
			State:      StateRelearning,
			Interval:   cfg.GraduatingInterval,
			EaseFactor: maxFloat(ease-cfg.EaseDecrement, cfg.MinEaseFactor),
			DueAt:      now.AddDate(0, 0, cfg.GraduatingInterval),
		}
	case RatingCorrect:
		newInterval := int(math.Round(float64(input.Interval) * ease))
		newInterval = maxInt(newInterval, input.Interval+1)
		return ScheduleOutput{
			State:      StateMastered,
			Interval:   newInterval,
			EaseFactor: ease,
			DueAt:      now.AddDate(0, 0, newInterval),
		}
	case RatingEasy:
		newInterval := int(math.Round(float64(input.Interval) * ease * cfg.EasyBonusMultipler))
		newInterval = maxInt(newInterval, input.Interval+1)
		return ScheduleOutput{
			State:      StateMastered,
			Interval:   newInterval,
			EaseFactor: minFloat(ease+cfg.EaseIncrement, cfg.MaxEaseFactor),
			DueAt:      now.AddDate(0, 0, newInterval),
		}
	default:
		return ScheduleOutput{
			State:      StateMastered,
			Interval:   input.Interval,
			EaseFactor: ease,
			DueAt:      now.AddDate(0, 0, input.Interval),
		}
	}
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

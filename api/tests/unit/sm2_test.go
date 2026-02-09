package unit

import (
	"testing"
	"time"

	"memwright/api/internal/srs"
)

func testSM2Config() srs.SM2Config {
	return srs.SM2Config{
		InitialEaseFactor:  2.5,
		MinEaseFactor:      1.3,
		MaxEaseFactor:      3.0,
		EaseDecrement:      0.2,
		EaseIncrement:      0.15,
		EasyBonusMultipler: 1.3,
		GraduatingInterval: 1,
		MasteredThreshold:  21,
	}
}

func TestSM2_Name(t *testing.T) {
	sm2 := srs.NewSM2(testSM2Config())
	if sm2.Name() != "sm2" {
		t.Errorf("expected name 'sm2', got '%s'", sm2.Name())
	}
}

func TestSM2_Config(t *testing.T) {
	cfg := srs.SM2Config{
		InitialEaseFactor:  2.0,
		MinEaseFactor:      1.5,
		MaxEaseFactor:      2.8,
		EaseDecrement:      0.1,
		EaseIncrement:      0.1,
		EasyBonusMultipler: 1.5,
		GraduatingInterval: 2,
		MasteredThreshold:  30,
	}
	sm2 := srs.NewSM2(cfg)

	if sm2.Config() != cfg {
		t.Errorf("expected config to match input")
	}
}

func TestSM2_NewCard_WrongRating(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateNew,
		EaseFactor: cfg.InitialEaseFactor,
	}

	output := sm2.Schedule(input, srs.RatingWrong, now)

	if output.State != srs.StateNew {
		t.Errorf("expected state %s, got %s", srs.StateNew, output.State)
	}
	if output.Interval != 0 {
		t.Errorf("expected interval 0, got %d", output.Interval)
	}
	expectedEase := cfg.InitialEaseFactor - cfg.EaseDecrement
	if output.EaseFactor != expectedEase {
		t.Errorf("expected ease factor %f, got %f", expectedEase, output.EaseFactor)
	}
	if !output.DueAt.Equal(now) {
		t.Errorf("expected due at %v, got %v", now, output.DueAt)
	}
}

func TestSM2_NewCard_CorrectRating(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateNew,
		EaseFactor: cfg.InitialEaseFactor,
	}

	output := sm2.Schedule(input, srs.RatingCorrect, now)

	if output.State != srs.StateLearning {
		t.Errorf("expected state %s, got %s", srs.StateLearning, output.State)
	}
	if output.Interval != cfg.GraduatingInterval {
		t.Errorf("expected interval %d, got %d", cfg.GraduatingInterval, output.Interval)
	}
	if output.EaseFactor != cfg.InitialEaseFactor {
		t.Errorf("expected ease factor %f, got %f", cfg.InitialEaseFactor, output.EaseFactor)
	}
	expectedDue := now.AddDate(0, 0, cfg.GraduatingInterval)
	if !output.DueAt.Equal(expectedDue) {
		t.Errorf("expected due at %v, got %v", expectedDue, output.DueAt)
	}
}

func TestSM2_NewCard_EasyRating(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateNew,
		EaseFactor: cfg.InitialEaseFactor,
	}

	output := sm2.Schedule(input, srs.RatingEasy, now)

	if output.State != srs.StateReview {
		t.Errorf("expected state %s, got %s", srs.StateReview, output.State)
	}
	expectedInterval := 3 // int(1 * 1.3 * 2.5) = 3
	if output.Interval != expectedInterval {
		t.Errorf("expected interval %d, got %d", expectedInterval, output.Interval)
	}
	expectedEase := cfg.InitialEaseFactor + cfg.EaseIncrement
	if output.EaseFactor != expectedEase {
		t.Errorf("expected ease factor %f, got %f", expectedEase, output.EaseFactor)
	}
}

func TestSM2_LearningCard_WrongRating(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateLearning,
		Interval:   1,
		EaseFactor: cfg.InitialEaseFactor,
	}

	output := sm2.Schedule(input, srs.RatingWrong, now)

	if output.State != srs.StateLearning {
		t.Errorf("expected state %s, got %s", srs.StateLearning, output.State)
	}
	if output.Interval != cfg.GraduatingInterval {
		t.Errorf("expected interval %d, got %d", cfg.GraduatingInterval, output.Interval)
	}
	expectedEase := cfg.InitialEaseFactor - cfg.EaseDecrement
	if output.EaseFactor != expectedEase {
		t.Errorf("expected ease factor %f, got %f", expectedEase, output.EaseFactor)
	}
}

func TestSM2_LearningCard_CorrectRating(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateLearning,
		Interval:   1,
		EaseFactor: cfg.InitialEaseFactor,
	}

	output := sm2.Schedule(input, srs.RatingCorrect, now)

	if output.State != srs.StateReview {
		t.Errorf("expected state %s, got %s", srs.StateReview, output.State)
	}
	if output.Interval < input.Interval+1 {
		t.Errorf("expected interval to increase, got %d", output.Interval)
	}
	if output.EaseFactor != cfg.InitialEaseFactor {
		t.Errorf("expected ease factor %f, got %f", cfg.InitialEaseFactor, output.EaseFactor)
	}
}

func TestSM2_ReviewCard_WrongRating_ResetsInterval(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateReview,
		Interval:   10,
		EaseFactor: cfg.InitialEaseFactor,
	}

	output := sm2.Schedule(input, srs.RatingWrong, now)

	if output.State != srs.StateRelearning {
		t.Errorf("expected state %s, got %s", srs.StateRelearning, output.State)
	}
	if output.Interval != cfg.GraduatingInterval {
		t.Errorf("expected interval to reset to %d, got %d", cfg.GraduatingInterval, output.Interval)
	}
	expectedEase := cfg.InitialEaseFactor - cfg.EaseDecrement
	if output.EaseFactor != expectedEase {
		t.Errorf("expected ease factor %f, got %f", expectedEase, output.EaseFactor)
	}
}

func TestSM2_ReviewCard_CorrectRating_MultipliesByEase(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateReview,
		Interval:   4,
		EaseFactor: 2.0,
	}

	output := sm2.Schedule(input, srs.RatingCorrect, now)

	if output.State != srs.StateReview {
		t.Errorf("expected state %s, got %s", srs.StateReview, output.State)
	}
	expectedInterval := 8 // 4 * 2.0
	if output.Interval != expectedInterval {
		t.Errorf("expected interval %d, got %d", expectedInterval, output.Interval)
	}
	if output.EaseFactor != 2.0 {
		t.Errorf("expected ease factor 2.0, got %f", output.EaseFactor)
	}
}

func TestSM2_ReviewCard_EasyRating_AppliesBonus(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateReview,
		Interval:   4,
		EaseFactor: 2.0,
	}

	output := sm2.Schedule(input, srs.RatingEasy, now)

	if output.State != srs.StateReview {
		t.Errorf("expected state %s, got %s", srs.StateReview, output.State)
	}
	expectedInterval := 10 // 4 * 2.0 * 1.3 = 10.4 rounded to 10
	if output.Interval != expectedInterval {
		t.Errorf("expected interval %d, got %d", expectedInterval, output.Interval)
	}
	expectedEase := 2.0 + cfg.EaseIncrement
	if output.EaseFactor != expectedEase {
		t.Errorf("expected ease factor %f, got %f", expectedEase, output.EaseFactor)
	}
}

func TestSM2_EaseFactorNeverDropsBelowMinimum(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateReview,
		Interval:   5,
		EaseFactor: cfg.MinEaseFactor,
	}

	output := sm2.Schedule(input, srs.RatingWrong, now)

	if output.EaseFactor != cfg.MinEaseFactor {
		t.Errorf("expected ease factor to stay at minimum %f, got %f", cfg.MinEaseFactor, output.EaseFactor)
	}
}

func TestSM2_EaseFactorNeverExceedsMaximum(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateReview,
		Interval:   5,
		EaseFactor: cfg.MaxEaseFactor,
	}

	output := sm2.Schedule(input, srs.RatingEasy, now)

	if output.EaseFactor > cfg.MaxEaseFactor {
		t.Errorf("expected ease factor to cap at %f, got %f", cfg.MaxEaseFactor, output.EaseFactor)
	}
}

func TestSM2_TransitionToMastered(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateReview,
		Interval:   10,
		EaseFactor: 2.5,
	}

	output := sm2.Schedule(input, srs.RatingCorrect, now)

	if output.Interval < cfg.MasteredThreshold {
		t.Fatalf("interval %d should be >= %d for this test", output.Interval, cfg.MasteredThreshold)
	}
	if output.State != srs.StateMastered {
		t.Errorf("expected state %s when interval >= %d, got %s", srs.StateMastered, cfg.MasteredThreshold, output.State)
	}
}

func TestSM2_MasteredCard_WrongRating_GoesToRelearning(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateMastered,
		Interval:   30,
		EaseFactor: cfg.InitialEaseFactor,
	}

	output := sm2.Schedule(input, srs.RatingWrong, now)

	if output.State != srs.StateRelearning {
		t.Errorf("expected state %s, got %s", srs.StateRelearning, output.State)
	}
	if output.Interval != cfg.GraduatingInterval {
		t.Errorf("expected interval to reset to %d, got %d", cfg.GraduatingInterval, output.Interval)
	}
}

func TestSM2_RelearningCard_CorrectRating(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateRelearning,
		Interval:   1,
		EaseFactor: 2.0,
	}

	output := sm2.Schedule(input, srs.RatingCorrect, now)

	if output.State != srs.StateReview {
		t.Errorf("expected state %s, got %s", srs.StateReview, output.State)
	}
	if output.Interval <= input.Interval {
		t.Errorf("expected interval to increase from %d, got %d", input.Interval, output.Interval)
	}
}

func TestSM2_DefaultEaseFactor_WhenZero(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateNew,
		EaseFactor: 0,
	}

	output := sm2.Schedule(input, srs.RatingCorrect, now)

	if output.EaseFactor != cfg.InitialEaseFactor {
		t.Errorf("expected ease factor %f when input is 0, got %f", cfg.InitialEaseFactor, output.EaseFactor)
	}
}

func TestSM2_DueAtCalculation(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateReview,
		Interval:   5,
		EaseFactor: 2.0,
	}

	output := sm2.Schedule(input, srs.RatingCorrect, now)

	expectedDue := now.AddDate(0, 0, output.Interval)
	if !output.DueAt.Equal(expectedDue) {
		t.Errorf("expected due at %v, got %v", expectedDue, output.DueAt)
	}
}

func TestSM2_IntervalAlwaysIncreases_OnCorrect(t *testing.T) {
	cfg := testSM2Config()
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		name     string
		interval int
		ease     float64
	}{
		{"small interval low ease", 1, cfg.MinEaseFactor},
		{"medium interval", 5, 2.0},
		{"large interval", 100, 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := srs.ScheduleInput{
				State:      srs.StateReview,
				Interval:   tc.interval,
				EaseFactor: tc.ease,
			}

			output := sm2.Schedule(input, srs.RatingCorrect, now)

			if output.Interval <= tc.interval {
				t.Errorf("expected interval to increase from %d, got %d", tc.interval, output.Interval)
			}
		})
	}
}

func TestSM2_ImplementsAlgorithmInterface(t *testing.T) {
	var _ srs.Algorithm = (*srs.SM2)(nil)
	var _ srs.Algorithm = srs.NewSM2(testSM2Config())
}

func TestSM2_CustomConfig_GraduatingInterval(t *testing.T) {
	cfg := srs.SM2Config{
		InitialEaseFactor:  2.5,
		MinEaseFactor:      1.3,
		MaxEaseFactor:      3.0,
		EaseDecrement:      0.2,
		EaseIncrement:      0.15,
		EasyBonusMultipler: 1.3,
		GraduatingInterval: 3,
		MasteredThreshold:  21,
	}
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateNew,
		EaseFactor: cfg.InitialEaseFactor,
	}

	output := sm2.Schedule(input, srs.RatingCorrect, now)

	if output.Interval != 3 {
		t.Errorf("expected custom graduating interval 3, got %d", output.Interval)
	}
	expectedDue := now.AddDate(0, 0, 3)
	if !output.DueAt.Equal(expectedDue) {
		t.Errorf("expected due at %v, got %v", expectedDue, output.DueAt)
	}
}

func TestSM2_CustomConfig_MasteredThreshold(t *testing.T) {
	cfg := srs.SM2Config{
		InitialEaseFactor:  2.5,
		MinEaseFactor:      1.3,
		MaxEaseFactor:      3.0,
		EaseDecrement:      0.2,
		EaseIncrement:      0.15,
		EasyBonusMultipler: 1.3,
		GraduatingInterval: 1,
		MasteredThreshold:  7,
	}
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateReview,
		Interval:   5,
		EaseFactor: 2.0,
	}

	output := sm2.Schedule(input, srs.RatingCorrect, now)

	if output.Interval < cfg.MasteredThreshold {
		t.Fatalf("expected interval >= %d, got %d", cfg.MasteredThreshold, output.Interval)
	}
	if output.State != srs.StateMastered {
		t.Errorf("expected state %s with custom threshold %d, got %s", srs.StateMastered, cfg.MasteredThreshold, output.State)
	}
}

func TestSM2_CustomConfig_EaseDecrement(t *testing.T) {
	cfg := srs.SM2Config{
		InitialEaseFactor:  2.5,
		MinEaseFactor:      1.3,
		MaxEaseFactor:      3.0,
		EaseDecrement:      0.3,
		EaseIncrement:      0.15,
		EasyBonusMultipler: 1.3,
		GraduatingInterval: 1,
		MasteredThreshold:  21,
	}
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateReview,
		Interval:   5,
		EaseFactor: 2.5,
	}

	output := sm2.Schedule(input, srs.RatingWrong, now)

	expectedEase := 2.5 - 0.3
	if output.EaseFactor != expectedEase {
		t.Errorf("expected ease factor %f with custom decrement, got %f", expectedEase, output.EaseFactor)
	}
}

func TestSM2_CustomConfig_EasyBonusMultiplier(t *testing.T) {
	cfg := srs.SM2Config{
		InitialEaseFactor:  2.5,
		MinEaseFactor:      1.3,
		MaxEaseFactor:      3.0,
		EaseDecrement:      0.2,
		EaseIncrement:      0.15,
		EasyBonusMultipler: 2.0,
		GraduatingInterval: 1,
		MasteredThreshold:  21,
	}
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateReview,
		Interval:   4,
		EaseFactor: 2.0,
	}

	output := sm2.Schedule(input, srs.RatingEasy, now)

	expectedInterval := 16 // 4 * 2.0 * 2.0 = 16
	if output.Interval != expectedInterval {
		t.Errorf("expected interval %d with custom bonus multiplier, got %d", expectedInterval, output.Interval)
	}
}

func TestSM2_CustomConfig_MinEaseFactor(t *testing.T) {
	cfg := srs.SM2Config{
		InitialEaseFactor:  2.5,
		MinEaseFactor:      1.5,
		MaxEaseFactor:      3.0,
		EaseDecrement:      0.2,
		EaseIncrement:      0.15,
		EasyBonusMultipler: 1.3,
		GraduatingInterval: 1,
		MasteredThreshold:  21,
	}
	sm2 := srs.NewSM2(cfg)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	input := srs.ScheduleInput{
		State:      srs.StateReview,
		Interval:   5,
		EaseFactor: 1.5,
	}

	output := sm2.Schedule(input, srs.RatingWrong, now)

	if output.EaseFactor != 1.5 {
		t.Errorf("expected ease factor to stay at custom minimum 1.5, got %f", output.EaseFactor)
	}
}

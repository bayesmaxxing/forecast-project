package models

import (
	"math"
	"testing"
	"time"
)

// Tests using resolvedAt as the close date
func TestCalcForecastScore_SinglePoint(t *testing.T) {
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC) // 7 days later

	points := []TimePoint{
		{PointForecast: 0.7, CreatedAt: forecastCreated},
	}

	// Test positive outcome (forecast resolves to YES)
	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, nil, &forecastResolved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// With single point, naive and time-weighted should be identical
	expectedBrier := math.Pow(0.7-1, 2) // (0.7-1)^2 = 0.09
	if math.Abs(score.BrierScore-expectedBrier) > 0.0001 {
		t.Errorf("BrierScore = %v, want %v", score.BrierScore, expectedBrier)
	}
	if math.Abs(score.BrierScoreTimeWeighted-expectedBrier) > 0.0001 {
		t.Errorf("BrierScoreTimeWeighted = %v, want %v", score.BrierScoreTimeWeighted, expectedBrier)
	}

	// Test negative outcome (forecast resolves to NO)
	score, err = CalcForecastScore(points, false, 1, 1, forecastCreated, nil, &forecastResolved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBrier = math.Pow(0.7, 2) // 0.7^2 = 0.49
	if math.Abs(score.BrierScore-expectedBrier) > 0.0001 {
		t.Errorf("BrierScore = %v, want %v", score.BrierScore, expectedBrier)
	}
}

func TestCalcForecastScore_MultiplePointsEqualDuration(t *testing.T) {
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC) // 3 days later

	// 3 points, each held for 1 day (equal duration)
	points := []TimePoint{
		{PointForecast: 0.3, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{PointForecast: 0.5, CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{PointForecast: 0.7, CreatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
	}

	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, nil, &forecastResolved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Naive average: (0.09 + 0.25 + 0.09) / 3 = 0.143333
	naiveBrier := (math.Pow(0.3-1, 2) + math.Pow(0.5-1, 2) + math.Pow(0.7-1, 2)) / 3
	if math.Abs(score.BrierScore-naiveBrier) > 0.0001 {
		t.Errorf("BrierScore = %v, want %v", score.BrierScore, naiveBrier)
	}

	// Time-weighted: each weight = 1/3, so same as naive
	// (0.09 * 1/3) + (0.25 * 1/3) + (0.09 * 1/3) = 0.143333
	expectedTimeWeighted := naiveBrier
	if math.Abs(score.BrierScoreTimeWeighted-expectedTimeWeighted) > 0.0001 {
		t.Errorf("BrierScoreTimeWeighted = %v, want %v", score.BrierScoreTimeWeighted, expectedTimeWeighted)
	}
}

func TestCalcForecastScore_MultiplePointsUnequalDuration(t *testing.T) {
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC) // 10 days later

	// Point 1: held for 1 day (10% of time)
	// Point 2: held for 9 days (90% of time)
	points := []TimePoint{
		{PointForecast: 0.9, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{PointForecast: 0.1, CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, nil, &forecastResolved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Naive: (0.01 + 0.81) / 2 = 0.41
	naiveBrier := (math.Pow(0.9-1, 2) + math.Pow(0.1-1, 2)) / 2
	if math.Abs(score.BrierScore-naiveBrier) > 0.0001 {
		t.Errorf("BrierScore = %v, want %v", score.BrierScore, naiveBrier)
	}

	// Time-weighted: (0.01 * 0.1) + (0.81 * 0.9) = 0.001 + 0.729 = 0.73
	expectedTimeWeighted := math.Pow(0.9-1, 2)*0.1 + math.Pow(0.1-1, 2)*0.9
	if math.Abs(score.BrierScoreTimeWeighted-expectedTimeWeighted) > 0.0001 {
		t.Errorf("BrierScoreTimeWeighted = %v, want %v (should heavily weight the 0.1 prediction held for 90%% of time)",
			score.BrierScoreTimeWeighted, expectedTimeWeighted)
	}

	// Time-weighted should be worse (higher) than naive here
	// because the bad prediction (0.1) was held much longer
	if score.BrierScoreTimeWeighted <= score.BrierScore {
		t.Errorf("Time-weighted score should be worse (higher) than naive when bad prediction held longer")
	}
}

func TestCalcForecastScore_WeightsSumToOne(t *testing.T) {
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)

	points := []TimePoint{
		{PointForecast: 0.5, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{PointForecast: 0.6, CreatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
		{PointForecast: 0.7, CreatedAt: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)},
	}

	// Manually calculate weights to verify they sum to 1.0
	totalOpenTime := forecastResolved.Sub(forecastCreated).Seconds()

	weight1 := points[1].CreatedAt.Sub(points[0].CreatedAt).Seconds() / totalOpenTime // 2 days
	weight2 := points[2].CreatedAt.Sub(points[1].CreatedAt).Seconds() / totalOpenTime // 5 days
	weight3 := forecastResolved.Sub(points[2].CreatedAt).Seconds() / totalOpenTime    // 3 days

	sumWeights := weight1 + weight2 + weight3
	if math.Abs(sumWeights-1.0) > 0.0001 {
		t.Errorf("Weights sum to %v, want 1.0", sumWeights)
	}
}

func TestCalcForecastScore_LogScores(t *testing.T) {
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	points := []TimePoint{
		{PointForecast: 0.8, CreatedAt: forecastCreated},
	}

	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, nil, &forecastResolved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// For outcome=true and prob=0.8:
	expectedLog2 := math.Log2(0.8)
	expectedLogN := math.Log(0.8)

	if math.Abs(score.Log2Score-expectedLog2) > 0.0001 {
		t.Errorf("Log2Score = %v, want %v", score.Log2Score, expectedLog2)
	}
	if math.Abs(score.LogNScore-expectedLogN) > 0.0001 {
		t.Errorf("LogNScore = %v, want %v", score.LogNScore, expectedLogN)
	}

	// Test negative outcome
	score, err = CalcForecastScore(points, false, 1, 1, forecastCreated, nil, &forecastResolved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// For outcome=false and prob=0.8:
	expectedLog2 = math.Log2(1 - 0.8)
	expectedLogN = math.Log(1 - 0.8)

	if math.Abs(score.Log2Score-expectedLog2) > 0.0001 {
		t.Errorf("Log2Score = %v, want %v", score.Log2Score, expectedLog2)
	}
	if math.Abs(score.LogNScore-expectedLogN) > 0.0001 {
		t.Errorf("LogNScore = %v, want %v", score.LogNScore, expectedLogN)
	}
}

func TestCalcForecastScore_EdgeCaseProbabilities(t *testing.T) {
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	// Test boundary values (0.0 and 1.0 should error)
	testCases := []struct {
		prob      float64
		shouldErr bool
	}{
		{0.0, true},    // exactly 0
		{1.0, true},    // exactly 1
		{0.001, false}, // just above 0
		{0.999, false}, // just below 1
		{-0.1, true},   // negative
		{1.1, true},    // above 1
	}

	for _, tc := range testCases {
		points := []TimePoint{
			{PointForecast: tc.prob, CreatedAt: forecastCreated},
		}

		_, err := CalcForecastScore(points, true, 1, 1, forecastCreated, nil, &forecastResolved)
		if tc.shouldErr && err == nil {
			t.Errorf("Expected error for probability %v, got nil", tc.prob)
		}
		if !tc.shouldErr && err != nil {
			t.Errorf("Unexpected error for probability %v: %v", tc.prob, err)
		}
	}
}

func TestCalcForecastScore_EmptyPoints(t *testing.T) {
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	points := []TimePoint{}

	_, err := CalcForecastScore(points, true, 1, 1, forecastCreated, nil, &forecastResolved)
	if err == nil {
		t.Error("Expected error for empty points slice, got nil")
	}
}

func TestCalcForecastScore_TimeWeightingFavorsLongerHeldPredictions(t *testing.T) {
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC) // 30 days

	// Someone who quickly corrected their prediction
	// Started at 0.2 (bad), held for 1 day
	// Then 0.9 (good), held for 29 days
	goodForecaster := []TimePoint{
		{PointForecast: 0.2, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{PointForecast: 0.9, CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	// Someone who made the opposite mistake
	// Started at 0.9 (good), held for 1 day
	// Then 0.2 (bad), held for 29 days
	badForecaster := []TimePoint{
		{PointForecast: 0.9, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{PointForecast: 0.2, CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	// Outcome = true (YES)
	goodScore, _ := CalcForecastScore(goodForecaster, true, 1, 1, forecastCreated, nil, &forecastResolved)
	badScore, _ := CalcForecastScore(badForecaster, true, 1, 1, forecastCreated, nil, &forecastResolved)

	// Naive scores should be identical (same probabilities, different order)
	if math.Abs(goodScore.BrierScore-badScore.BrierScore) > 0.0001 {
		t.Errorf("Naive scores should be identical, got %v vs %v", goodScore.BrierScore, badScore.BrierScore)
	}

	// Time-weighted: good forecaster should have much better (lower) score
	if goodScore.BrierScoreTimeWeighted >= badScore.BrierScoreTimeWeighted {
		t.Errorf("Good forecaster should have better time-weighted score. Got good=%v, bad=%v",
			goodScore.BrierScoreTimeWeighted, badScore.BrierScoreTimeWeighted)
	}

	// The difference should be substantial (at least 0.3 difference in Brier)
	diff := badScore.BrierScoreTimeWeighted - goodScore.BrierScoreTimeWeighted
	if diff < 0.3 {
		t.Errorf("Expected substantial difference in time-weighted scores, got %v", diff)
	}
}

func TestCalcForecastScore_MetadataCorrect(t *testing.T) {
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	points := []TimePoint{
		{PointForecast: 0.5, CreatedAt: forecastCreated},
	}

	var userID int64 = 42
	var forecastID int64 = 123

	score, err := CalcForecastScore(points, true, userID, forecastID, forecastCreated, nil, &forecastResolved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if score.UserID != userID {
		t.Errorf("UserID = %v, want %v", score.UserID, userID)
	}
	if score.ForecastID != forecastID {
		t.Errorf("ForecastID = %v, want %v", score.ForecastID, forecastID)
	}
	if score.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

// Tests for closeDate calculation logic

func TestCalcForecastScore_CloseDateNil(t *testing.T) {
	// When forecastClosingDate is nil, should use forecastResolvedAt as closeDate
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC) // 10 days later

	// Point 1: held for 5 days (50% of time from creation to resolve)
	// Point 2: held for 5 days (50% of time from creation to resolve)
	points := []TimePoint{
		{PointForecast: 0.9, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{PointForecast: 0.1, CreatedAt: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)},
	}

	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, nil, &forecastResolved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// When closeDate is nil, it uses forecastResolvedAt
	// Point 1: 5 days / 10 days = 0.5
	// Point 2: 5 days / 10 days = 0.5
	expectedTimeWeighted := math.Pow(0.9-1, 2)*0.5 + math.Pow(0.1-1, 2)*0.5
	if math.Abs(score.BrierScoreTimeWeighted-expectedTimeWeighted) > 0.0001 {
		t.Errorf("BrierScoreTimeWeighted = %v, want %v (should use forecastResolvedAt as closeDate when closingDate is nil)",
			score.BrierScoreTimeWeighted, expectedTimeWeighted)
	}
}

func TestCalcForecastScore_CloseDateBeforeResolved(t *testing.T) {
	// When forecastClosingDate is before forecastResolvedAt, should use forecastClosingDate
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastClosing := time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)   // 5 days after creation
	forecastResolved := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC) // 10 days after creation, 5 days after closing

	// Point 1: held for 2 days (40% of the 5-day closing window)
	// Point 2: held for 3 days (60% of the 5-day closing window)
	points := []TimePoint{
		{PointForecast: 0.9, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{PointForecast: 0.1, CreatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
	}

	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, &forecastClosing, &forecastResolved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Time weights based on closing date (5 days total)
	// Point 1: 2 days / 5 days = 0.4
	// Point 2: 3 days / 5 days = 0.6
	expectedTimeWeighted := math.Pow(0.9-1, 2)*0.4 + math.Pow(0.1-1, 2)*0.6
	if math.Abs(score.BrierScoreTimeWeighted-expectedTimeWeighted) > 0.0001 {
		t.Errorf("BrierScoreTimeWeighted = %v, want %v (should use forecastClosingDate as closeDate)",
			score.BrierScoreTimeWeighted, expectedTimeWeighted)
	}
}

func TestCalcForecastScore_CloseDateAfterResolved(t *testing.T) {
	// When forecastClosingDate is after forecastResolvedAt, should use forecastResolvedAt as closeDate
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC) // 5 days after creation
	forecastClosing := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC) // 10 days after creation, after resolved

	points := []TimePoint{
		{PointForecast: 0.7, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{PointForecast: 0.3, CreatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
	}

	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, &forecastClosing, &forecastResolved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// When closeDate=forecastClosing is after resolution, should fall back to forecastResolvedAt
	// Point 1: 2 days / 5 days = 0.4
	// Point 2: 3 days / 5 days = 0.6
	expectedTimeWeighted := math.Pow(0.7-1, 2)*0.4 + math.Pow(0.3-1, 2)*0.6
	if math.Abs(score.BrierScoreTimeWeighted-expectedTimeWeighted) > 0.0001 {
		t.Errorf("BrierScoreTimeWeighted = %v, want %v (should use forecastResolvedAt when closingDate is after resolvedAt)",
			score.BrierScoreTimeWeighted, expectedTimeWeighted)
	}
}

func TestCalcForecastScore_CloseDateEqualsResolved(t *testing.T) {
	// When forecastClosingDate equals forecastResolvedAt (not Before), should use forecastResolvedAt
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)
	forecastClosing := forecastResolved // Same time

	points := []TimePoint{
		{PointForecast: 0.8, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{PointForecast: 0.4, CreatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
	}

	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, &forecastClosing, &forecastResolved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Since closingDate.Before(resolvedAt) is false when they're equal,
	// should use forecastResolvedAt
	// Point 1: 2 days / 5 days = 0.4
	// Point 2: 3 days / 5 days = 0.6
	expectedTimeWeighted := math.Pow(0.8-1, 2)*0.4 + math.Pow(0.4-1, 2)*0.6
	if math.Abs(score.BrierScoreTimeWeighted-expectedTimeWeighted) > 0.0001 {
		t.Errorf("BrierScoreTimeWeighted = %v, want %v (should use forecastResolvedAt when closingDate equals resolvedAt)",
			score.BrierScoreTimeWeighted, expectedTimeWeighted)
	}
}

func TestCalcForecastScore_AllPointsBeforeCloseDate(t *testing.T) {
	// All forecast points are made before the closing date
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastClosing := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)  // 10 days after creation
	forecastResolved := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // 14 days after creation

	// All points before closing date with UNEQUAL durations
	// Point 1: Days 1-3 (2 days = 20%)
	// Point 2: Days 3-11 (8 days = 80%)
	points := []TimePoint{
		{PointForecast: 0.2, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{PointForecast: 0.8, CreatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
	}

	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, &forecastClosing, &forecastResolved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Time-weighted based on 10-day closing window
	// Point 1: 2 days / 10 days = 0.2
	// Point 2: 8 days / 10 days = 0.8
	expectedTimeWeighted := math.Pow(0.2-1, 2)*0.2 + math.Pow(0.8-1, 2)*0.8
	if math.Abs(score.BrierScoreTimeWeighted-expectedTimeWeighted) > 0.0001 {
		t.Errorf("BrierScoreTimeWeighted = %v, want %v",
			score.BrierScoreTimeWeighted, expectedTimeWeighted)
	}

	// Verify that using closing date affects the score differently than using resolved date
	scoreWithoutClosing, _ := CalcForecastScore(points, true, 1, 1, forecastCreated, nil, &forecastResolved)

	// When closeDate is nil, it uses forecastResolvedAt (14 days), resulting in:
	// Point 1: 2 days / 14 days ≈ 0.143
	// Point 2: 12 days / 14 days ≈ 0.857
	// This is different from using closingDate (10 days):
	// Point 1: 2 days / 10 days = 0.2
	// Point 2: 8 days / 10 days = 0.8
	// These should be different
	if math.Abs(score.BrierScoreTimeWeighted-scoreWithoutClosing.BrierScoreTimeWeighted) < 0.0001 {
		t.Error("Closing date should affect time-weighted scores differently than no closing date")
	}
}

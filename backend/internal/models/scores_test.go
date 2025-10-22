package models

import (
	"math"
	"testing"
	"time"
)

func TestCalcForecastScore_SinglePoint(t *testing.T) {
	forecastCreated := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	forecastResolved := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC) // 7 days later

	points := []TimePoint{
		{PointForecast: 0.7, CreatedAt: forecastCreated},
	}

	// Test positive outcome (forecast resolves to YES)
	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, &forecastResolved)
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
	score, err = CalcForecastScore(points, false, 1, 1, forecastCreated, &forecastResolved)
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

	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, &forecastResolved)
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

	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, &forecastResolved)
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
	weight3 := forecastResolved.Sub(points[2].CreatedAt).Seconds() / totalOpenTime     // 3 days

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

	score, err := CalcForecastScore(points, true, 1, 1, forecastCreated, &forecastResolved)
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
	score, err = CalcForecastScore(points, false, 1, 1, forecastCreated, &forecastResolved)
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
		{0.0, true},   // exactly 0
		{1.0, true},   // exactly 1
		{0.001, false}, // just above 0
		{0.999, false}, // just below 1
		{-0.1, true},  // negative
		{1.1, true},   // above 1
	}

	for _, tc := range testCases {
		points := []TimePoint{
			{PointForecast: tc.prob, CreatedAt: forecastCreated},
		}

		_, err := CalcForecastScore(points, true, 1, 1, forecastCreated, &forecastResolved)
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

	_, err := CalcForecastScore(points, true, 1, 1, forecastCreated, &forecastResolved)
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
	goodScore, _ := CalcForecastScore(goodForecaster, true, 1, 1, forecastCreated, &forecastResolved)
	badScore, _ := CalcForecastScore(badForecaster, true, 1, 1, forecastCreated, &forecastResolved)

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

	score, err := CalcForecastScore(points, true, userID, forecastID, forecastCreated, &forecastResolved)
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

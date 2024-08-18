package utils

import (
	"errors"
	"math"
)

type ForecastScores struct {
	BrierScore float64
	Log2Score  float64
	LogNScore  float64
}

func CalcForecastScores(probabilities []float64, outcome bool) (ForecastScores, error) {
	if len(probabilities) == 0 {
		return ForecastScores{}, errors.New("no probabilities provided")
	}

	var brierSum, log2Sum, logNSum float64
	points := float64(len(probabilities))

	for _, prob := range probabilities {
		if prob <= 0 || prob >= 1 {
			return ForecastScores{}, errors.New("probs must be within 0 and 1")
		}

		if outcome {
			brierSum += math.Pow(prob-1, 2)
			logNSum += math.Log(prob)
			log2Sum += math.Log2(prob)
		} else {
			brierSum += math.Pow(prob, 2)
			logNSum += math.Log(1 - prob)
			log2Sum += math.Log2(1 - prob)
		}
	}

	return ForecastScores{
		BrierScore: brierSum / points,
		Log2Score:  log2Sum / points,
		LogNScore:  logNSum / points,
	}, nil
}

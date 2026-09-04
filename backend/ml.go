package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/malaschitz/randomForest"
)

// TrainRFModel reads a CSV containing the points data, trains a Random Forest classifier,
// evaluates it on the test set (West Kameng district), and returns the feature importance weights.
func TrainRFModel(csvPath string) (map[string]float64, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV is empty or has only header")
	}

	header := records[0]
	// Expected features
	featureNames := []string{"slope", "twi", "curvature", "landcover_class", "elevation", "dist_to_drainage"}
	featureIndices := make([]int, len(featureNames))
	for i, name := range featureNames {
		idx := -1
		for j, h := range header {
			if strings.TrimSpace(h) == name {
				idx = j
				break
			}
		}
		if idx == -1 {
			return nil, fmt.Errorf("feature %s not found in CSV header", name)
		}
		featureIndices[i] = idx
	}

	districtIdx := -1
	targetIdx := -1
	for j, h := range header {
		if strings.TrimSpace(h) == "district" {
			districtIdx = j
		}
		if strings.TrimSpace(h) == "is_landslide" {
			targetIdx = j
		}
	}

	if districtIdx == -1 {
		return nil, fmt.Errorf("district column not found")
	}
	if targetIdx == -1 {
		return nil, fmt.Errorf("is_landslide column not found")
	}

	var trainX, testX [][]float64
	var trainY, testY []int

	for i := 1; i < len(records); i++ {
		record := records[i]
		district := record[districtIdx]

		x := make([]float64, len(featureNames))
		for j, fIdx := range featureIndices {
			val, _ := strconv.ParseFloat(record[fIdx], 64)
			x[j] = val
		}

		target, _ := strconv.Atoi(record[targetIdx])

		if district == "West Kameng" {
			testX = append(testX, x)
			testY = append(testY, target)
		} else {
			trainX = append(trainX, x)
			trainY = append(trainY, target)
		}
	}

	if len(trainX) == 0 {
		return nil, fmt.Errorf("no training data found")
	}

	fmt.Printf("Training Random Forest with %d samples... (Test samples: %d)\n", len(trainX), len(testX))

	forest := &randomforest.Forest{
		Data: randomforest.ForestData{
			X:     trainX,
			Class: trainY,
		},
	}
	// Train 300 trees
	forest.Train(300)

	// Calculate AUC if test data exists
	if len(testX) > 0 {
		auc := calculateAUC(forest, testX, testY)
		fmt.Printf("AUC: %.4f\n", auc)
	}

	// Extract feature importances and normalize them
	importances := forest.FeatureImportance
	sum := 0.0
	for _, imp := range importances {
		sum += imp
	}

	weights := make(map[string]float64)
	for i, name := range featureNames {
		if i < len(importances) {
			if sum > 0 {
				weights[name] = importances[i] / sum
			} else {
				weights[name] = 0.0
			}
		}
	}

	return weights, nil
}

type predResult struct {
	prob float64
	label int
}

func calculateAUC(forest *randomforest.Forest, X [][]float64, Y []int) float64 {
	var preds []predResult
	posCount := 0
	negCount := 0

	for i, x := range X {
		// Vote returns votes for each class
		votes := forest.Vote(x)
		// Assuming binary classification (0 and 1)
		prob1 := 0.0
		if len(votes) > 1 {
			prob1 = votes[1] // We can just use the vote count for class 1
		}
		
		preds = append(preds, predResult{prob: prob1, label: Y[i]})
		if Y[i] == 1 {
			posCount++
		} else {
			negCount++
		}
	}

	if posCount == 0 || negCount == 0 {
		return 0.0 // AUC not defined if only one class is present
	}

	// Sort predictions by probability descending
	sort.Slice(preds, func(i, j int) bool {
		return preds[i].prob > preds[j].prob
	})

	// Wilcoxon-Mann-Whitney statistic
	sumRanks := 0.0
	for i, p := range preds {
		if p.label == 1 {
			// Rank is len(preds) - i
			sumRanks += float64(len(preds) - i)
		}
	}

	auc := (sumRanks - float64(posCount*(posCount+1))/2.0) / float64(posCount*negCount)
	return auc
}

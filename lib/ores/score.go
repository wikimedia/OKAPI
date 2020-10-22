package ores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"okapi/lib/cache"
	"strconv"
	"strings"
)

// Available ORES models
const (
	Damaging Model = "damaging"
)

// Model ORES assestment model
type Model string

// CanScore check if database can be scored by model
func (model Model) CanScore(dbName string) bool {
	return cache.Client().SIsMember(getCacheKey(dbName), string(model)).Val()
}

// ScoreOne score one by predefined model
func (model Model) ScoreOne(dbName string, revision int) (*Score, error) {
	return ScoreOne(dbName, revision, model)
}

// ScoreMany score many revisions by one model
func (model Model) ScoreMany(dbName string, revisions []int) (map[int]*Score, error) {
	return ScoreMany(dbName, model, revisions)
}

// Probability model probability
type Probability struct {
	True  float64 `json:"true"`
	False float64 `json:"false"`
}

// Score model score
type Score struct {
	Prediction  bool        `json:"prediction"`
	Probability Probability `json:"probability"`
}

// Scores model scores
type Scores map[int]map[Model]Score

// Stream stream scores
type Stream struct {
	ModelName   Model `json:"model_name"`
	Probability Probability
}

// ScoreOne score one model by revision
func ScoreOne(dbName string, revision int, model Model) (*Score, error) {
	if !cache.Client().SIsMember(getCacheKey(dbName), string(model)).Val() {
		return nil, fmt.Errorf("Error: no '%s' model for '%s' database", model, dbName)
	}

	res, err := client.R().Get("scores/" + dbName + "/" + strconv.Itoa(revision) + "/" + string(model))

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Error: status code '%d'", res.StatusCode())
	}

	oresScores := map[string]struct {
		Score Score `json:"score"`
	}{}

	if err = json.Unmarshal(res.Body(), &oresScores); err != nil {
		return nil, err
	}

	if modelScore, exists := oresScores[string(model)]; exists {
		return &modelScore.Score, nil
	}

	return nil, fmt.Errorf("Error: no '%s' model for '%s' database", model, dbName)
}

// ScoreMany score many revisions by one model
func ScoreMany(dbName string, model Model, revisions []int) (map[int]*Score, error) {
	if !cache.Client().SIsMember(getCacheKey(dbName), string(model)).Val() {
		return nil, fmt.Errorf("Error: no '%s' model for '%s' database", model, dbName)
	}

	if len(revisions) <= 0 {
		return nil, fmt.Errorf("Error: empty revisions list")
	}

	params := ""

	for _, revision := range revisions {
		params += "rev_id=" + strconv.Itoa(revision) + "&"
	}

	res, err := client.R().Get("scores/" + dbName + "/" + string(model) + "?" + strings.TrimRight(params, "&"))

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Error: status code '%d'", res.StatusCode())
	}

	oresScores := map[int]map[string]struct {
		Score Score `json:"score"`
	}{}

	if err = json.Unmarshal(res.Body(), &oresScores); err != nil {
		return nil, err
	}

	scores := map[int]*Score{}

	for rev, modelScores := range oresScores {
		if modelScore, exists := modelScores[string(model)]; exists {
			scores[rev] = &modelScore.Score
		}
	}

	return scores, nil
}

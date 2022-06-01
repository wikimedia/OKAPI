package schema

import ores "github.com/protsack-stephan/mediawiki-ores-client"

// Scores ORES scores representation
type Scores struct {
	Damaging  *ores.ScoreDamaging  `json:"damaging,omitempty"`
	GoodFaith *ores.ScoreGoodFaith `json:"goodfaith,omitempty"`
}

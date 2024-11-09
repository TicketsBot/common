package model

import "github.com/rxdn/gdl/objects/guild"

type MisconductAlert struct {
	Guild      *guild.Guild   `json:"guild"`
	Score      int            `json:"score"`
	RuleScores map[string]int `json:"rule_scores"`
}

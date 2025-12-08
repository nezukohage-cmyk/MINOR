package models

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type QuizAttempt struct {
	mgm.DefaultModel `bson:",inline"`

	UserID       string         `json:"user_id" bson:"user_id"`
	QuizSession  string         `json:"quiz_session" bson:"quiz_session"`
	Total        int            `json:"total" bson:"total"`
	Correct      int            `json:"correct" bson:"correct"`
	Accuracy     float64        `json:"accuracy" bson:"accuracy"`
	TimeTakenSec int            `json:"time_taken_sec" bson:"time_taken_sec"`
	ByTopic      map[string]any `json:"by_topic" bson:"by_topic"` // e.g. wrong counts per topic
	CreatedAt    time.Time      `json:"created_at" bson:"created_at"`
}

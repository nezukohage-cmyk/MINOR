package models

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type QuizSession struct {
	mgm.DefaultModel `bson:",inline"`

	UserID         string         `json:"user_id" bson:"user_id"`
	QuestionIDs    []string       `json:"question_ids" bson:"question_ids"`
	RequestedCount map[string]int `json:"requested_count" bson:"requested_count"`
	Meta           map[string]any `json:"meta" bson:"meta"`
	StartedAt      time.Time      `json:"started_at" bson:"started_at"`
	SubmittedAt    time.Time      `json:"submitted_at,omitempty" bson:"submitted_at,omitempty"`
	TimeTakenSec   int            `json:"time_taken_sec,omitempty" bson:"time_taken_sec,omitempty"`
	Score          int            `json:"score" bson:"score"`
}

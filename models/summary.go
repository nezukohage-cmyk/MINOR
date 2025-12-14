package models

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type Summary struct {
	mgm.DefaultModel `bson:",inline"`

	UserID    string    `json:"user_id" bson:"user_id"`
	InputText string    `json:"input_text" bson:"input_text"`
	Summary   string    `json:"summary" bson:"summary"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

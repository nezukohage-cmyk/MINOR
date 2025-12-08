package models

import "github.com/kamva/mgm/v3"

type TodoTask struct {
	mgm.DefaultModel `bson:",inline"`

	UserID string `json:"user_id" bson:"user_id"`
	PlanID string `json:"plan_id" bson:"plan_id"`
	Date   string `json:"date" bson:"date"`
	Task   string `json:"task" bson:"task"`
	Done   bool   `json:"done" bson:"done"`
}

package models

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type DailyTask struct {
	Date  string   `json:"date" bson:"date"`
	Tasks []string `json:"tasks" bson:"tasks"`
}

type StudyPlan struct {
	mgm.DefaultModel `bson:",inline"`

	UserID      string                 `json:"user_id" bson:"user_id"`
	Title       string                 `json:"title" bson:"title"`
	InputRows   []StudyPlannerRow      `json:"input_rows" bson:"input_rows"`
	PlanJSON    map[string]interface{} `json:"plan_json" bson:"plan_json"`
	RawOutput   string                 `json:"raw_output" bson:"raw_output"`
	DailyTable  []DailyTask            `json:"daily_timetable" bson:"daily_timetable"`
	GeneratedAt time.Time              `json:"generated_at" bson:"generated_at"`
}

// StudyPlannerRow = one row from the user's planner UI
type StudyPlannerRow struct {
	Subject     string   `json:"subject" bson:"subject"`
	Topics      []string `json:"topics" bson:"topics"`
	Deadline    string   `json:"deadline" bson:"deadline"` // ISO date string or user text
	HoursPerDay int      `json:"hours_per_day" bson:"hours_per_day"`
	HoursNeeded int      `json:"hours_needed" bson:"hours_needed"` // optional, ai may change
}

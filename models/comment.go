package models

import "github.com/kamva/mgm/v3"

type Comment struct {
	mgm.DefaultModel `bson:",inline"`

	NoteID   string `json:"note_id" bson:"note_id"`
	UserID   string `json:"user_id" bson:"user_id"`
	Text     string `json:"text" bson:"text"`
	ParentID string `json:"parent_id,omitempty" bson:"parent_id,omitempty"`

	Upvotes   []string `json:"upvotes" bson:"upvotes"`
	Downvotes []string `json:"downvotes" bson:"downvotes"`
	Score     int      `json:"score" bson:"score"`
}

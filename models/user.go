package models

import "github.com/kamva/mgm/v3"

type User struct {
	mgm.DefaultModel `bson:",inline"`

	Username     string   `json:"username" bson:"username"`
	Email        string   `json:"email,omitempty" bson:"email,omitempty"`
	Phone        string   `json:"phone,omitempty" bson:"phone,omitempty"`
	Password     string   `json:"password" bson:"password"`
	CommunityIDs []string `json:"community_ids" bson:"community_ids"`
	SavedNoteIDs []string `json:"saved_note_ids" bson:"saved_note_ids"`
	Interests    []string `json:"interests" bson:"interests"`
}

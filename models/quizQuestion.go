package models

import "github.com/kamva/mgm/v3"

type Question struct {
	mgm.DefaultModel `bson:",inline"`

	Text      string   `json:"text" bson:"text"`             // question text
	Options   []string `json:"options" bson:"options"`       // 4 options
	Answer    string   `json:"answer" bson:"answer"`         // correct option
	SubjectID string   `json:"subject_id" bson:"subject_id"` // subject tag id or name
	TopicIDs  []string `json:"topic_ids" bson:"topic_ids"`   // multiple topics
	CreatedBy string   `json:"created_by" bson:"created_by"` // userID
}

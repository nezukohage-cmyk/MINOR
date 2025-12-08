package models

import "github.com/kamva/mgm/v3"

type Tag struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string   `json:"name" bson:"name"`
	Type             string   `json:"type" bson:"type"`
	ChildTopics      []string `json:"child_topics" bson:"child_topics"`
	ParentSubjects   []string `json:"parent_subjects" bson:"parent_subjects"`
}

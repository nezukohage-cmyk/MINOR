package models

import "github.com/kamva/mgm/v3"

type Note struct {
	mgm.DefaultModel `bson:",inline"`

	UserID           string   `json:"user_id" bson:"user_id"`
	SubjectIDs       []string `json:"subject_ids" bson:"subject_ids"`
	TopicIDs         []string `json:"topic_ids" bson:"topic_ids"`
	FileName         string   `json:"file_name" bson:"file_name"`
	FileType         string   `json:"file_type" bson:"file_type"`
	Url              string   `json:"url" bson:"url"`
	CloudinaryID     string   `json:"cloudinary_id" bson:"cloudinary_id"`
	Size             int64    `json:"size" bson:"size"`
	ModerationStatus string   `json:"moderation_status" bson:"moderation_status"`

	Upvotes   []string `json:"upvotes" bson:"upvotes"`
	Downvotes []string `json:"downvotes" bson:"downvotes"`
	Score     int      `json:"score" bson:"score"`

	CommunityID string   `json:"community_id" bson:"community_id"`
	Urls        []string `json:"urls" bson:"urls"`
	PublicIDs   []string `json:"public_ids" bson:"public_ids"`
}

type Subject struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
}

type Topic struct {
	mgm.DefaultModel `bson:",inline"`
	Subject          string `json:"subject" bson:"subject"`
	Name             string `json:"name" bson:"name"`
}

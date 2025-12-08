// package models

// import "github.com/kamva/mgm/v3"

// type Note struct {
// 	mgm.DefaultModel `bson:",inline"`

//		UserID           string `json:"user_id" bson:"user_id"`
//		SubjectID        string `json:"subject_id" bson:"subject_id"`
//		TopicID          string `json:"topic_id" bson:"topic_id"`
//		FileName         string `json:"file_name" bson:"file_name"`
//		FileType         string `json:"file_type" bson:"file_type"`
//		Url              string `json:"url" bson:"url"`
//		CloudinaryID     string `json:"cloudinary_id" bson:"cloudinary_id"`
//		Size             int64  `json:"size" bson:"size"`
//		ModerationStatus string `json:"moderation_status" bson:"moderation_status"`
//	}
package models

import "github.com/kamva/mgm/v3"

type Note struct {
	mgm.DefaultModel `bson:",inline"` // <-- important

	UserID           string   `json:"user_id" bson:"user_id"`
	SubjectIDs       []string `json:"subject_ids" bson:"subject_ids"`
	TopicIDs         []string `json:"topic_ids" bson:"topic_ids"`
	FileName         string   `json:"file_name" bson:"file_name"`
	FileType         string   `json:"file_type" bson:"file_type"`
	Url              string   `json:"url" bson:"url"`
	CloudinaryID     string   `json:"cloudinary_id" bson:"cloudinary_id"`
	Size             int64    `json:"size" bson:"size"`
	ModerationStatus string   `json:"moderation_status" bson:"moderation_status"`

	// Voting
	Upvotes   []string `json:"upvotes" bson:"upvotes"`
	Downvotes []string `json:"downvotes" bson:"downvotes"`
	Score     int      `json:"score" bson:"score"`

	// Community
	CommunityID   string `json:"community_id" bson:"community_id"`
	CommunityName string `json:"community_name" bson:"community_name"`
}

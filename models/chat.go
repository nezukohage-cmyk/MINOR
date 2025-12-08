package models

import "github.com/kamva/mgm/v3"

type ChatMessage struct {
	Role string `json:"role" bson:"role"` // user / assistant
	Text string `json:"text" bson:"text"`
}

type Chat struct {
	mgm.DefaultModel `bson:",inline"`

	UserID   string        `json:"user_id" bson:"user_id"`
	Title    string        `json:"title" bson:"title"`
	Messages []ChatMessage `json:"messages" bson:"messages"`
}

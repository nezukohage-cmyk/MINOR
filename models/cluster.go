package models

import "github.com/kamva/mgm/v3"

type Cluster struct {
	mgm.DefaultModel `bson:",inline"`

	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	Tags        []string `json:"tags" bson:"tags"`
}

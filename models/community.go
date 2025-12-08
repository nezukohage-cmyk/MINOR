package models

import "github.com/kamva/mgm/v3"

type Community struct {
	mgm.DefaultModel `bson:",inline"`

	Name      string   `json:"name" bson:"name"`     // e.g. "SDMCET Dharwad"
	Region    string   `json:"region" bson:"region"` // e.g. "Dharwad"
	Members   []string `json:"members" bson:"members"`
	CreatedBy string   `json:"created_by" bson:"created_by"`
}

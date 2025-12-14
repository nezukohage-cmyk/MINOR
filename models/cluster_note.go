package models

import "github.com/kamva/mgm/v3"

type ClusterNote struct {
	mgm.DefaultModel `bson:",inline"`

	ClusterID  string `json:"cluster_id" bson:"cluster_id"`
	UploaderID string `json:"uploader_id" bson:"uploader_id"`

	Title string   `json:"title" bson:"title"`
	Tags  []string `json:"tags" bson:"tags"`

	FileURL string `json:"file_url" bson:"file_url"`
}

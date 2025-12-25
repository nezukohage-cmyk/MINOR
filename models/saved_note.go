package models

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SavedNote struct {
	mgm.DefaultModel `bson:",inline"`

	UserID     primitive.ObjectID `bson:"user_id"`
	NoteID     primitive.ObjectID `bson:"note_id"`
	ClusterID  string             `bson:"cluster_id"`
	UploaderID string             `bson:"uploader_id"`
	SavedAt    time.Time          `bson:"saved_at"`
	Title      string             `bson:"title"`
	Tags       []string           `bson:"tags"`
	FileURL    string             `bson:"file_url"`
}

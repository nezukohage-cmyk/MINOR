package controllers

import (
	"net/http"
	"time"

	"lexxi/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveNote(c *gin.Context) {
	userIDStr := c.GetString("user_id")

	var body struct {
		NoteID    string   `json:"note_id" binding:"required"`
		ClusterID string   `json:"cluster_id" binding:"required"`
		Title     string   `json:"title" binding:"required"`
		FileURL   string   `json:"file_url" binding:"required"`
		Tags      []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	noteID, err := primitive.ObjectIDFromHex(body.NoteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	// Prevent duplicate saves
	filter := bson.M{
		"user_id": userID,
		"note_id": noteID,
	}

	count, _ := mgm.Coll(&models.SavedNote{}).CountDocuments(c, filter)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "already saved"})
		return
	}

	saved := &models.SavedNote{
		UserID:    userID,
		NoteID:    noteID,
		ClusterID: body.ClusterID,

		Title:   body.Title,
		FileURL: body.FileURL,
		Tags:    body.Tags,

		SavedAt: time.Now(),
	}

	if err := mgm.Coll(saved).Create(saved); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}

func GetSavedNotes(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	var saved []models.SavedNote
	err := mgm.Coll(&models.SavedNote{}).
		SimpleFind(&saved, bson.M{"user_id": userID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch saved notes"})
		return
	}

	var result []gin.H
	for _, s := range saved {
		result = append(result, gin.H{
			"id":       s.NoteID.Hex(),
			"title":    s.Title,
			"file_url": s.FileURL,
			"tags":     s.Tags,
			"saved_at": s.SavedAt,
		})

	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
func UnsaveNote(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	var body struct {
		NoteID string `json:"note_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	noteID, err := primitive.ObjectIDFromHex(body.NoteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	_, err = mgm.Coll(&models.SavedNote{}).DeleteOne(
		c,
		bson.M{
			"user_id": userID,
			"note_id": noteID,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unsave"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unsaved"})
}

package controllers

import (
	"net/http"
	"sort"

	//"strings"
	"time"

	"lexxi/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func HomeFeed(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	// 1️⃣ Fetch all approved notes
	var notes []models.Note
	err := mgm.Coll(&models.Note{}).SimpleFind(
		&notes,
		bson.M{"moderation_status": "approved"},
	)
	if err != nil {
		notes = []models.Note{}
	}

	// Sort newest → oldest
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].CreatedAt.After(notes[j].CreatedAt)
	})

	// Limit to 10 items
	if len(notes) > 10 {
		notes = notes[:10]
	}

	// 2️⃣ Build feed list
	type feedItem struct {
		CreatedAt time.Time              `json:"created_at"`
		Type      string                 `json:"type"`
		Data      map[string]interface{} `json:"data"`
	}

	items := []feedItem{}

	for _, n := range notes {
		items = append(items, feedItem{
			CreatedAt: n.CreatedAt,
			Type:      "note",
			Data: map[string]interface{}{
				"id":            n.ID.Hex(),
				"user_id":       n.UserID,
				"subject_ids":   n.SubjectIDs,
				"topic_ids":     n.TopicIDs,
				"file_name":     n.FileName,
				"file_type":     n.FileType,
				"url":           n.Url,
				"cloudinary_id": n.CloudinaryID,
				"size":          n.Size,
				"score":         n.Score,
			},
		})
	}

	// 3️⃣ SAFETY: ensure response is never null
	if items == nil {
		items = []feedItem{}
	}

	c.JSON(http.StatusOK, gin.H{"feed": items})
}

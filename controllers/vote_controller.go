package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"

	"lexxi/models"
)

func Vote(c *gin.Context, voteType string) {
	userID := c.GetString("user_id")
	noteID := c.Param("id")

	note := &models.Note{}
	if err := mgm.Coll(note).FindByID(noteID, note); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	// Remove existing vote from both lists
	note.Upvotes = removeUser(note.Upvotes, userID)
	note.Downvotes = removeUser(note.Downvotes, userID)

	// Apply new vote if user toggled differently
	switch voteType {
	case "up":
		note.Upvotes = append(note.Upvotes, userID)
	case "down":
		note.Downvotes = append(note.Downvotes, userID)
	}

	// Recalculate score
	note.Score = len(note.Upvotes) - len(note.Downvotes)

	if err := mgm.Coll(note).Update(note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"score":     note.Score,
		"upvoted":   contains(note.Upvotes, userID),
		"downvoted": contains(note.Downvotes, userID),
	})
}

func Upvote(c *gin.Context) {
	Vote(c, "up")
}

func Downvote(c *gin.Context) {
	Vote(c, "down")
}

func removeUser(list []string, id string) []string {
	newList := []string{}
	for _, v := range list {
		if v != id {
			newList = append(newList, v)
		}
	}
	return newList
}

func contains(list []string, id string) bool {
	for _, v := range list {
		if v == id {
			return true
		}
	}
	return false
}

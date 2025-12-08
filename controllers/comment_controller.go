package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"

	"lexxi/models"
)

func CreateComment(c *gin.Context) {
	userID := c.GetString("user_id")
	noteID := c.Param("id")

	var body struct {
		Text     string `json:"text"`
		ParentID string `json:"parent_id"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	comment := &models.Comment{
		NoteID:    noteID,
		UserID:    userID,
		Text:      body.Text,
		ParentID:  body.ParentID,
		Upvotes:   []string{},
		Downvotes: []string{},
		Score:     0,
	}

	if err := mgm.Coll(comment).Create(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "comment added", "id": comment.ID.Hex()})
}
func GetComments(c *gin.Context) {
	noteID := c.Param("id")

	var comments []models.Comment
	err := mgm.Coll(&models.Comment{}).SimpleFind(&comments, bson.M{"note_id": noteID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comments"})
		return
	}

	// Build nested tree
	tree := buildCommentTree(comments)

	c.JSON(http.StatusOK, tree)
}

func buildCommentTree(comments []models.Comment) []gin.H {
	idMap := map[string]gin.H{}
	tree := []gin.H{}

	for _, c := range comments {
		node := gin.H{
			"id":      c.ID.Hex(),
			"user_id": c.UserID,
			"text":    c.Text,
			"score":   c.Score,
			"replies": []gin.H{},
		}
		idMap[c.ID.Hex()] = node
	}

	for _, c := range comments {
		if c.ParentID == "" {
			tree = append(tree, idMap[c.ID.Hex()])
		} else {
			parent, ok := idMap[c.ParentID]
			if ok {
				parent["replies"] = append(parent["replies"].([]gin.H), idMap[c.ID.Hex()])
			}
		}
	}

	return tree
}

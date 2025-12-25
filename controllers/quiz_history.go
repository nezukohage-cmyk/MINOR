package controllers

import (
	"net/http"

	"lexxi/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func GetQuizHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var sessions []models.QuizSession
	err := mgm.Coll(&models.QuizSession{}).
		SimpleFind(&sessions, bson.M{"user_id": userID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load history"})
		return
	}

	out := make([]gin.H, 0, len(sessions))

	for _, s := range sessions {
		out = append(out, gin.H{
			"quiz_id":          s.ID.Hex(),
			"subjects":         s.Subjects,
			"score":            s.Score,
			"served_questions": len(s.QuestionIDs),
			"started_at":       s.StartedAt,
			"submitted_at":     s.SubmittedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": out})
}

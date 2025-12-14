package controllers

import (
	"lexxi/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SummarizeText(c *gin.Context) {

	var body struct {
		Text string `json:"text" binding:"required"`
	}

	// Validate request
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
		return
	}

	if strings.TrimSpace(body.Text) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "text cannot be empty"})
		return
	}

	summary, err := services.Summarize(body.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
	})
}
func GetSummaries(c *gin.Context) {
	userID := c.GetString("user_id")

	// Fetch summaries from DB
	items, err := services.FetchSummaries(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch summaries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"summaries": items})
}

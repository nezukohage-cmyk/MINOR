package controllers

import (
	"lexxi/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AskChatbot(c *gin.Context) {

	var body struct {
		Subject  string `json:"subject"`  // optional
		Question string `json:"question"` // required
	}

	// Parse JSON
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Question is mandatory
	if strings.TrimSpace(body.Question) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "question is required"})
		return
	}

	// Call Gemini
	answer, err := services.AskGemini(body.Subject, body.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"subject":  body.Subject,
		"question": body.Question,
		"answer":   answer,
	})
}

package controllers

import (
	"net/http"
	"strings"

	"lexxi/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// helper → normalize string (UPPERCASE + TRIM)
func norm(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

// helper → normalize topic list
func normalizeTopicList(list []string) []string {
	out := []string{}
	for _, t := range list {
		nt := norm(t)
		if nt != "" {
			out = append(out, nt)
		}
	}
	return out
}

func CreateQuestion(c *gin.Context) {
	userID := c.GetString("user_id")

	// INPUTS
	text := strings.TrimSpace(c.PostForm("text"))
	subject := norm(c.PostForm("subject")) // normalize
	topics := normalizeTopicList(c.PostFormArray("topics"))

	optionInputs := c.PostFormArray("options") // should be 4
	answer := norm(c.PostForm("answer"))       // normalize

	// VALIDATION
	if text == "" || subject == "" || len(topics) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing fields"})
		return
	}

	// CHECK DUPLICATE QUESTION (case-insensitive)
	filter := bson.M{
		"text": bson.M{"$regex": "^" + text + "$", "$options": "i"},
	}

	var existing models.Question
	err := mgm.Coll(&models.Question{}).First(filter, &existing)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "duplicate question"})
		return
	}

	// VALIDATE OPTIONS
	if len(optionInputs) != 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exactly 4 options required"})
		return
	}

	normalizedOptions := []string{}
	optionSet := make(map[string]bool)

	for _, opt := range optionInputs {
		nopt := norm(opt)

		if nopt == "" || optionSet[nopt] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "options must be unique and non-empty"})
			return
		}

		optionSet[nopt] = true
		normalizedOptions = append(normalizedOptions, nopt)
	}

	// ANSWER MUST MATCH ONE OPTION
	if !optionSet[answer] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "answer must match one of the options"})
		return
	}

	// CREATE QUESTION OBJECT
	q := &models.Question{
		Text:      text,
		Options:   normalizedOptions,
		Answer:    answer,
		SubjectID: subject,
		TopicIDs:  topics,
		CreatedBy: userID,
	}

	// SAVE
	if err := mgm.Coll(q).Create(q); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save question"})
		return
	}

	// SUCCESS
	c.JSON(http.StatusCreated, gin.H{
		"message": "question added",
		"id":      q.ID.Hex(),
	})
}
